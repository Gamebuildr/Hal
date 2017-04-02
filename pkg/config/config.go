package config

// Auth0ClientSecret is the secret key for authenticating tokens
const Auth0ClientSecret string = "AUTH0_CLIENT_SECRET"

// LogEndpoint is the endpoint for sending log data
const LogEndpoint string = "PAPERTRAIL_ENDPOINT"

// Region is the world region messages are coming from
const Region string = "REGION"

// QueueURL is the URL endpoint to recieve messages
const QueueURL string = "QUEUE_URL"

// MrrobotNotifications is the URL endpoint to send MrRobot messages
const MrrobotNotifications string = "MRROBOT_NOTIFICATIONS"

// GamebuildrNotifications is the URL endpoint to send Gamebuildr messages
const GamebuildrNotifications string = "GAMEBUILDR_NOTIFICATIONS"

// CodeRepoStorage is the location to save source code
const CodeRepoStorage string = "CODE_REPO_STORAGE"

// GoEnv is the environment the current system is operating in
const GoEnv string = "GO_ENV"

// AWSAccessKeyId is the key id for using amazon services
const AWSAccessKeyId string = "AWS_ACCESS_KEY_ID"

// AWSAccessKey is the secret key for using amazon services
const AWSAccessKey string = "AWS_SECRET_ACCESS_KEY"

// BuildSourcePath is the source code for game engines
const BuildSourcePath string = "BUILD_SOURCE_PATH"

// BuildTargetPath is the target path for game builds
const BuildTargetPath string = "BUILD_TARGET_PATH"

// EngineLogPath is the log path output from game engines
const EngineLogPath string = "ENGINE_LOG_PATH"
