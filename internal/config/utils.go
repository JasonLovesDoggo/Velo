package config

// ToEnv converts the Environment map to a slice of strings in KEY=VALUE format
func (s *ServiceDefinition) ToEnv() []string {
	env := make([]string, 0, len(s.Environment))
	for key, value := range s.Environment {
		env = append(env, key+"="+value)
	}
	return env
}
