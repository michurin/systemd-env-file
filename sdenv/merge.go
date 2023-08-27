package sdenv

import "os"

func Environ(env []string, filenames ...string) ([]string, error) {
	x := append([]string(nil), env...) // make local copy to avoid side effects in case F(x[:1])
	for _, name := range filenames {
		cfgData, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		pairs, err := Parser([]rune(string(cfgData)))
		if err != nil {
			return nil, err
		}
		for _, v := range pairs {
			x = append(x, v[0]+"="+v[1])
		}
	}
	return x, nil
}
