package darksky

// https://openweathermap.org/api/one-call-api
type owmForecastResponse struct {
	Latitude       Measurement   `json:"lat"`
	Longitude      Measurement   `json:"lon"`
	Timezone       string        `json:"timezone"`
	TimezoneOffset float64       `json:"timezone_offset"`
	Current        owmCurrent    `json:"current,omitempty"`
	Minutely       []owmMinutely `json:"minutely,omitempty"`
	Hourly         []owmHourly   `json:"hourly,omitempty"`
	Daily          []owmDaily    `json:"daily,omitempty"`
}

type owmWeather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type owmMinutely struct {
	DataTime      Timestamp   `json:"dt"`
	Precipitation Measurement `json:"precipitation"`
}

type owmCurrent struct {
	DataTime    Timestamp    `json:"dt"`
	Sunrise     Timestamp    `json:"sunrise"`
	Sunset      Timestamp    `json:"sunset"`
	Temperature Measurement  `json:"temp"`
	FeelsLike   Measurement  `json:"feels_like"`
	Pressure    Measurement  `json:"pressure"`
	Humidity    Measurement  `json:"humidity"`
	DewPoint    Measurement  `json:"dew_point"`
	UvIndex     Measurement  `json:"uvi"`
	Clouds      Measurement  `json:"clouds"`
	Visibility  Measurement  `json:"visibility"`
	WindSpeed   Measurement  `json:"wind_speed"`
	WindDeg     Measurement  `json:"wind_deg"`
	WindGust    Measurement  `json:"wind_gust"`
	Weather     []owmWeather `json:"weather,omitempty"`
}

type owmHourly struct {
	DataTime    Timestamp    `json:"dt"`
	Temperature Measurement  `json:"temp"`
	FeelsLike   Measurement  `json:"feels_like"`
	Pressure    Measurement  `json:"pressure"`
	Humidity    Measurement  `json:"humidity"`
	DewPoint    Measurement  `json:"dew_point"`
	UvIndex     Measurement  `json:"uvi"`
	Clouds      Measurement  `json:"clouds"`
	Visibility  Measurement  `json:"visibility"`
	WindSpeed   Measurement  `json:"wind_speed"`
	WindDeg     Measurement  `json:"wind_deg"`
	WindGust    Measurement  `json:"wind_gust"`
	Weather     []owmWeather `json:"weather,omitempty"`
	Pop         Measurement  `json:"pop"`
	Rain        struct {
		OneH Measurement `json:"1h"`
	} `json:"rain,omitempty"`
	Snow struct {
		OneH Measurement `json:"1h"`
	} `json:"snow,omitempty"`
}

type owmDaily struct {
	DataTime    Timestamp `json:"dt"`
	Sunrise     Timestamp `json:"sunrise"`
	Sunset      Timestamp `json:"sunset"`
	Moonrise    Timestamp `json:"moonrise"`
	Moonset     Timestamp `json:"moonset"`
	MoonPhase   float64   `json:"moon_phase"`
	Temperature struct {
		Day   Measurement `json:"day"`
		Min   Measurement `json:"min"`
		Max   Measurement `json:"max"`
		Night Measurement `json:"night"`
		Eve   Measurement `json:"eve"`
		Morn  Measurement `json:"morn"`
	} `json:"temp"`
	FeelsLike struct {
		Day   Measurement `json:"day"`
		Night Measurement `json:"night"`
		Eve   Measurement `json:"eve"`
		Morn  Measurement `json:"morn"`
	} `json:"feels_like"`
	Pressure   Measurement  `json:"pressure"`
	Humidity   Measurement  `json:"humidity"`
	DewPoint   Measurement  `json:"dew_point"`
	WindSpeed  Measurement  `json:"wind_speed"`
	WindDeg    Measurement  `json:"wind_deg"`
	WindGust   Measurement  `json:"wind_gust"`
	Weather    []owmWeather `json:"weather,omitempty"`
	Clouds     Measurement  `json:"clouds"`
	Visibility Measurement  `json:"visibility"`
	Pop        Measurement  `json:"pop"`
	Rain       Measurement  `json:"rain,omitempty"`
	Snow       Measurement  `json:"snow,omitempty"`
	UvIndex    Measurement  `json:"uvi"`
}
