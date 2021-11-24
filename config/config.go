package config

type Properties struct {
	Port            string `env:"APP_PORT" env-default:"3000"`
	Host            string `env:"HOST" env-default:"localhost"`
	MailPort        int    `env:"MAIL_PORT" env-default:"587"`
	MailHost        string `env:"MAIL_HOST" env-default:"smtp.gmail.com"`
	MailFromAddress string `env:"MAIL_FROM" env-default:"bwrepairmanagement@gmail.com"`
	MailPassword    string `env:"MAIL_PASSWORD" env-default:"Imeidb@123"`
}
