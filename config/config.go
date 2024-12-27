package config

import (
	"embed"
	"flag"
	"fmt"
	"github.com/lyh-demo/go-webapp-demo/util"
	"gopkg.in/yaml.v3"
	"os"
)

type DatabaseConfig struct {
	Dialect   string `default:"sqlite3"`
	Host      string `default:"book.db"`
	Port      string
	Dbname    string
	Username  string
	Password  string
	Migration bool `default:"false"`
}
type RedisConfig struct {
	Enabled            bool `default:"false"`
	ConnectionPoolSize int  `yaml:"connection_pool_size" default:"10"`
	Host               string
	Port               string
}
type ExtensionConfig struct {
	MasterGenerator bool `yaml:"master_generator" default:"false"`
	CorsEnabled     bool `yaml:"cors_enabled" default:"false"`
	SecurityEnabled bool `yaml:"security_enabled" default:"false"`
}
type LogConfig struct {
	RequestLogFormat string `yaml:"request_log_format" default:"${remote_ip} ${account_name} ${uri} ${method} ${status}"`
}
type StaticContentsConfig struct {
	Enabled bool `default:"false"`
}
type SwaggerConfig struct {
	Enabled bool `default:"false"`
	Path    string
}
type SecurityConfig struct {
	AuthPath    []string `yaml:"auth_path"`
	ExcludePath []string `yaml:"exclude_path"`
	UserPath    []string `yaml:"user_path"`
	AdminPath   []string `yaml:"admin_path"`
}

// Config represents the composition of yml settings.
type Config struct {
	Database       DatabaseConfig       `yaml:"database"`
	Redis          RedisConfig          `yaml:"redis"`
	Extension      ExtensionConfig      `yaml:"extension"`
	Log            LogConfig            `yaml:"log"`
	StaticContents StaticContentsConfig `yaml:"static_contents"`
	Swagger        SwaggerConfig        `yaml:"swagger"`
	Security       SecurityConfig       `yaml:"security"`
}

const (
	// DEV represents development environment
	DEV = "development"
	// PRD represents production environment
	PRD = "production"
	// DOC represents docker container
	DOC = "docker"
)

// LoadAppConfig reads the settings written to the yml file
func LoadAppConfig(yamlFile embed.FS) (*Config, string) {
	var env *string
	if value := os.Getenv("WEB_APP_ENV"); value != "" {
		env = &value
	} else {
		env = flag.String("env", "development", "To switch configurations.")
		flag.Parse()
	}

	file, err := yamlFile.ReadFile(fmt.Sprintf(AppConfigPath, *env))
	if err != nil {
		fmt.Printf("Failed to read application.%s.yml: %s", *env, err)
		os.Exit(ErrExitStatus)
	}

	config := &Config{}
	if err := yaml.Unmarshal(file, config); err != nil {
		fmt.Printf("Failed to read application.%s.yml: %s", *env, err)
		os.Exit(ErrExitStatus)
	}

	return config, *env
}

// LoadMessagesConfig loads the messages.properties.
func LoadMessagesConfig(propsFile embed.FS) map[string]string {
	messages := util.ReadPropertiesFile(propsFile, MessagesConfigPath)
	if messages == nil {
		fmt.Printf("Failed to load the messages.properties.")
		os.Exit(ErrExitStatus)
	}
	return messages
}
