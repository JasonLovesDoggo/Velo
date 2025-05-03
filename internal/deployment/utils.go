package deployment

func (s *ServiceDefinition) ToEnv() []string {
	var env []string
	for k, v := range s.Environment {
		env = append(env, k+"="+v)
	}
	return env
}
