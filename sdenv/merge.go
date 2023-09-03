package sdenv

import "os"

// Environ merges environments in os.Environ()-format with
// variables sourced from given files and returns results
// in the same format, that is suitable for exec.Command().Env
// for instance.
func Environ(env []string, filenames ...string) ([]string, error) {
	x := append([]string(nil), env...) // make local copy to avoid side effects in case F(x[:1])
	for _, name := range filenames {
		cfgData, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		pairs, err := Parser(cfgData)
		if err != nil {
			return nil, err
		}
		for _, v := range pairs {
			x = append(x, v[0]+"="+v[1])
		}
	}
	return x, nil
}
