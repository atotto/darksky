package darksky

import (
	"gonum.org/v1/gonum/interp"
)

// https://openweathermap.org/darksky-openweather
func convert(o *owmForecastResponse, r *ForecastResponse, prediction string) error {
	r.Latitude = o.Latitude
	r.Longitude = o.Longitude
	r.Timezone = o.Timezone
	r.Offset = o.TimezoneOffset / 3600
	r.Currently = convertCurrent(&o.Current)
	r.Minutely = convertMinutely(o.Minutely)
	r.Hourly = convertHourly(o.Hourly, o.Daily, prediction)
	r.Daily = convertDaily(o.Daily)
	return nil
}

func convertCurrent(o *owmCurrent) *DataPoint {
	var summary, icon string
	if len(o.Weather) > 0 {
		w := o.Weather[0]
		summary = w.Description
		icon = w.Icon
	}
	return &DataPoint{
		Time:                 o.DataTime,
		Summary:              summary,
		Icon:                 icon,
		NearestStormDistance: 0,
		NearestStormBearing:  0,
		PrecipIntensity:      0,
		PrecipIntensityError: 0,
		PrecipProbability:    0,
		Temperature:          o.Temperature,
		ApparentTemperature:  o.FeelsLike,
		DewPoint:             o.DewPoint,
		Humidity:             o.Humidity,
		Pressure:             o.Pressure,
		WindSpeed:            o.WindSpeed,
		WindGust:             o.WindGust,
		WindBearing:          o.WindDeg,
		CloudCover:           o.Clouds / 100,
		UvIndex:              int64(o.UvIndex),
		Visibility:           o.Visibility,
		Ozone:                0,

		SunriseTime: o.Sunrise,
		SunsetTime:  o.Sunset,
	}
}

func convertMinutely(minutely []owmMinutely) *DataBlock {
	data := make([]DataPoint, 0, len(minutely))
	for _, o := range minutely {
		data = append(data, DataPoint{
			Time:            o.DataTime,
			PrecipIntensity: o.Precipitation,
		})
	}

	return &DataBlock{
		Data:    data,
		Summary: "",
		Icon:    "",
	}
}

func convertHourly(hourly []owmHourly, daily []owmDaily, prediction string) *DataBlock {
	var timestamp Timestamp
	data := make([]DataPoint, 0, len(hourly)+(len(daily)-len(hourly)/24)*24)
	for _, o := range hourly {
		var summary, icon string
		if len(o.Weather) > 0 {
			w := o.Weather[0]
			summary = w.Description
			icon = w.Icon
		}
		data = append(data, DataPoint{
			Time:                o.DataTime,
			Summary:             summary,
			Icon:                icon,
			PrecipIntensity:     o.Rain.OneH + o.Snow.OneH,
			PrecipProbability:   o.Pop,
			Temperature:         o.Temperature,
			ApparentTemperature: o.FeelsLike,
			DewPoint:            o.DewPoint,
			Humidity:            o.Humidity,
			Pressure:            o.Pressure,
			WindSpeed:           o.WindSpeed,
			WindGust:            o.WindGust,
			WindBearing:         o.WindDeg,
			CloudCover:          o.Clouds / 100,
			UvIndex:             int64(o.UvIndex),
			Visibility:          o.Visibility,
			Ozone:               0,
		})
		timestamp = o.DataTime
	}

	if prediction == "" {
		return &DataBlock{
			Data:    data,
			Summary: "",
			Icon:    "",
		}
	}

	size := (len(daily) - len(hourly)/24) * 24
	temperatures := make([]float64, 0, size)
	timestamps := make([]float64, 0, size)

	cloudCovers := make([]float64, 0, size)
	precipIntensities := make([]float64, 0, size)
	precipProbabilities := make([]float64, 0, size)
	dayTimestamps := make([]float64, 0, size)

	last := hourly[len(hourly)-1]
	cloudCovers = append(cloudCovers, float64(last.Clouds)/100)
	precipIntensities = append(precipIntensities, float64(last.Rain.OneH+last.Snow.OneH))
	precipProbabilities = append(precipProbabilities, float64(last.Pop))
	dayTimestamps = append(dayTimestamps, float64(last.DataTime))

	for _, d := range daily[len(hourly)/24:] {
		temperatures = append(temperatures, float64(d.Temperature.Morn), float64(d.Temperature.Day), float64(d.Temperature.Eve), float64(d.Temperature.Night))
		timestamps = append(timestamps, float64(d.Sunrise), float64(d.Sunrise+d.Sunset)/2, float64(d.Sunset), float64(d.Sunset)+3600*3)

		if d.DataTime > last.DataTime {
			cloudCovers = append(cloudCovers, float64(d.Clouds)/100)
			dayTimestamps = append(dayTimestamps, float64(d.Sunrise+d.Sunset)/2)
			precipIntensities = append(precipIntensities, float64(d.Rain+d.Snow)/24)
			precipProbabilities = append(precipProbabilities, float64(d.Pop))
		}
	}

	var pT, pC, pP, pPP interp.FittablePredictor
	switch prediction {
	case "linear":
		pT = &interp.PiecewiseLinear{}
		pC = &interp.PiecewiseLinear{}
		pP = &interp.PiecewiseLinear{}
		pPP = &interp.PiecewiseLinear{}
	case "constant":
		pT = &interp.PiecewiseConstant{}
		pC = &interp.PiecewiseConstant{}
		pP = &interp.PiecewiseConstant{}
		pPP = &interp.PiecewiseConstant{}
	case "spline":
		pT = &interp.AkimaSpline{}
		pC = &interp.AkimaSpline{}
		pP = &interp.AkimaSpline{}
		pPP = &interp.AkimaSpline{}
	default:
	}

	if err := pT.Fit(timestamps, temperatures); err != nil {
		return nil
	}
	if err := pC.Fit(dayTimestamps, cloudCovers); err != nil {
		return nil
	}
	if err := pP.Fit(dayTimestamps, precipIntensities); err != nil {
		return nil
	}
	if err := pPP.Fit(dayTimestamps, precipProbabilities); err != nil {
		return nil
	}

	for _, _ = range daily[len(hourly)/24+1:] {
		for i := 0; i < 24; i++ {
			timestamp += 3600
			data = append(data, DataPoint{
				Time:              timestamp,
				Temperature:       Measurement(pT.Predict(float64(timestamp))),
				CloudCover:        Measurement(pC.Predict(float64(timestamp))),
				PrecipIntensity:   Measurement(pP.Predict(float64(timestamp))),
				PrecipProbability: Measurement(pPP.Predict(float64(timestamp))),
			})
		}
	}

	return &DataBlock{
		Data:    data,
		Summary: "",
		Icon:    "",
	}
}

func convertDaily(daily []owmDaily) *DataBlock {
	data := make([]DataPoint, 0, len(daily))
	for _, o := range daily {
		var summary, icon string
		if len(o.Weather) > 0 {
			w := o.Weather[0]
			summary = w.Description
			icon = w.Icon
		}
		data = append(data, DataPoint{
			Time:        o.DataTime,
			Summary:     summary,
			Icon:        icon,
			SunriseTime: o.Sunrise,
			SunsetTime:  o.Sunset,
			MoonPhase:   0,

			PrecipIntensity:        o.Rain + o.Snow,
			PrecipIntensityMax:     0,
			PrecipIntensityMaxTime: 0,
			PrecipProbability:      0,
			PrecipType:             "",

			TemperatureHigh:     o.Temperature.Day,
			TemperatureHighTime: 0,
			TemperatureLow:      o.Temperature.Night,
			TemperatureLowTime:  0,

			ApparentTemperatureHigh:     o.FeelsLike.Day,
			ApparentTemperatureHighTime: 0,
			ApparentTemperatureLow:      o.FeelsLike.Night,
			ApparentTemperatureLowTime:  0,

			DewPoint:     o.DewPoint,
			Humidity:     o.Humidity,
			Pressure:     o.Pressure,
			WindSpeed:    o.WindSpeed,
			WindGust:     o.WindGust,
			WindBearing:  o.WindDeg,
			WindGustTime: 0,
			CloudCover:   o.Clouds / 100,
			UvIndex:      int64(o.UvIndex),
			UvIndexTime:  0,
			Visibility:   o.Visibility,
			Ozone:        0,

			TemperatureMin:     o.Temperature.Min,
			TemperatureMinTime: 0,
			TemperatureMax:     o.Temperature.Max,
			TemperatureMaxTime: 0,

			ApparentTemperatureMin:     0,
			ApparentTemperatureMinTime: 0,
			ApparentTemperatureMax:     0,
			ApparentTemperatureMaxTime: 0,
		})
	}

	return &DataBlock{
		Data:    data,
		Summary: "",
		Icon:    "",
	}
}
