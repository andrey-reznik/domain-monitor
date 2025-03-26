package configuration

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfiguration struct {
	// The port the application listens on
	Port int `yaml:"port" json:"port" default:"3124"`
	// Allow automatic WHOIS refresh
	AutomateWHOISRefresh bool `yaml:"automateWHOISRefresh" json:"automateWHOISRefresh" default:"true"`
	// Show the configuration in the web interface. This is a security risk and should be disabled in production
	ShowConfiguration bool `yaml:"showConfiguration" json:"showConfiguration" default:"false"`
}

type AlertsConfiguration struct {
	// The admin email address for receiving alerts
	Admin string `yaml:"admin" json:"admin"`
	// The admin telegram chatID for receiving alerts
	TelegramAdmin string `yaml:"telegramAdmin" json:"telegramAdmin"`
	// Send alerts for monitored domains
	SendAlerts bool `yaml:"sendAlerts" json:"sendAlerts"`
	// Send 2-month alert for domain expiry date
	Send2MonthAlert bool `yaml:"send2MonthAlert" json:"send2MonthAlert"`
	// Send 1-month alert for domain expiry date
	Send1MonthAlert bool `yaml:"send1MonthAlert" json:"send1MonthAlert" default:"true"`
	// Send 2-week alert for domain expiry date
	Send2WeekAlert bool `yaml:"send2WeekAlert" json:"send2WeekAlert"`
	// Send 1-week alert for domain expiry date
	Send1WeekAlert bool `yaml:"send1WeekAlert" json:"send1WeekAlert"`
	// Send 3-day alert for domain expiry date
	Send3DayAlert bool `yaml:"send3DayAlert" json:"send3DayAlert" default:"true"`
	// Send daily alerts within 7 days of domain expiry
	SendDailyExpiryAlert bool `yaml:"sendDailyExpiryAlert" json:"sendDailyExpiryAlert"`
}

type SMTPConfiguration struct {
	// SMTP host
	Host string `yaml:"host" json:"host"`
	// SMTP port
	Port int `yaml:"port" json:"port"`
	// Use secure connection (TLS)
	Secure bool `yaml:"secure" json:"secure"`
	// SMTP user name
	AuthUser string `yaml:"authUser" json:"authUser"`
	// SMTP user password
	AuthPass string `yaml:"authPass" json:"authPass"`
	// Enable SMTP
	Enabled bool `yaml:"enabled" json:"enabled"`
	// Name of the sender
	FromName string `yaml:"fromName" json:"fromName"`
	// Email address of the sender
	FromAddress string `yaml:"fromAddress" json:"fromAddress"`
}

type SchedulerConfiguration struct {
	// Interval after which WHOIS cache data is considered stale (in days)
	WhoisCacheStaleInterval int `yaml:"whoisCacheStaleInterval" json:"whoisCacheStaleInterval"`
	// Use standard WHOIS refresh schedule:
	//
	// 0. Cache miss for domain
	// 1. Cache becomes WhoisCacheStaleInterval old
	// 2. 3 months before expiry
	// 3. 2 months before expiry
	// 4. 1 month before expiry
	// 5. 2 weeks before expiry
	//
	// As always, manual refresh is possible, and can be triggered via the API or the web interface
	UseStandardWhoisRefreshSchedule bool `yaml:"useStandardWhoisRefreshSchedule" json:"useStandardWhoisRefreshSchedule"`
}

type TelegramConfiguration struct {
	// Telegram bot id
	BotID string `yaml:"botId" json:"botId"`
	// Enable Telegram
	Enabled bool `yaml:"enabled" json:"enabled"`
}

type ConfigurationFile struct {
	// The application configuration
	App AppConfiguration `yaml:"app" json:"app"`
	// The alerts configuration
	Alerts AlertsConfiguration `yaml:"alerts" json:"alerts"`
	// The SMTP configuration
	SMTP SMTPConfiguration `yaml:"smtp" json:"smtp"`
	// The scheduler configuration
	Scheduler SchedulerConfiguration `yaml:"scheduler" json:"scheduler"`
	// The Telegram configuration
	Telegram TelegramConfiguration `yaml:"telegram" json:"telegram"`
}

type Configuration struct {
	// The config data
	Config ConfigurationFile
	// The path to the config file
	Filepath string
}

// returns default configuration
func DefaultConfiguration(filepath string) Configuration {
	return Configuration{
		Filepath: filepath,
		Config: ConfigurationFile{
			App: AppConfiguration{
				Port:                 3124,
				AutomateWHOISRefresh: true,
				ShowConfiguration:    true,
			},
			Scheduler: SchedulerConfiguration{
				WhoisCacheStaleInterval:         190,
				UseStandardWhoisRefreshSchedule: true,
			},
			Alerts: AlertsConfiguration{
				Send1MonthAlert: true,
				Send3DayAlert:   true,
			},
		},
	}
}

// Write the app configuration to the config file
func (c Configuration) Flush() {
	data, dataErr := yaml.Marshal(c.Config)
	if dataErr != nil {
		log.Println("Error while marshalling configuration")
		log.Fatalf("error: %v", dataErr)
	}

	file, err := os.Create(c.Filepath)
	if err != nil {
		log.Println("Error while creating configuration file")
		log.Fatalf("error: %v", err)
	}

	defer file.Close()

	_, err = io.WriteString(file, string(data))
	if err != nil {
		log.Println("Error while writing configuration file")
		log.Fatalf("error: %v", err)
	}

	// Check if the file has been written
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("Error while checking configuration file")
		log.Fatalf("error: %v", err)
	}

	log.Printf("💾 Configuration flushed to %s", fileInfo.Name())
}

// Update the app configuration with the given data
func (c *Configuration) UpdateAppConfiguration(data AppConfiguration) {
	c.Config.App = data

	c.Flush()
}

// Update the alerts configuration with the given data
func (c *Configuration) UpdateAlertsConfiguration(data AlertsConfiguration) {
	c.Config.Alerts = data

	c.Flush()
}

// Update the SMTP configuration with the given data
func (c *Configuration) UpdateSMTPConfiguration(data SMTPConfiguration) {
	c.Config.SMTP = data

	c.Flush()
}

// Update the Telegram configuration with the given data
func (c *Configuration) UpdateTelegramConfiguration(data TelegramConfiguration) {
	c.Config.Telegram = data

	c.Flush()
}

// Update the scheduler configuration with the given data
func (c *Configuration) UpdateSchedulerConfiguration(data SchedulerConfiguration) {
	c.Config.Scheduler = data

	c.Flush()
}
