package config

type errorConf string

func (err errorConf) Error() string { return string(err) }

const errorEnvNotFound = errorConf("environment variable was not found")
const errorEnvIncorrectFormat = errorConf("incorrect format of the environment variable")
