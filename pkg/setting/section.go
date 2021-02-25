package setting

type RedisSettingS struct {
	Host string
	Name string
	Port string
}

type ServerSettingS struct {
	Port string
}


func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
