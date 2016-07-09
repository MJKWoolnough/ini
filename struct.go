package ini

import "reflect"

func (d *decoder) readStruct(s reflect.Value) {
	t := s.Type()
	nf := t.NumField()
Loop:
	for d.Peek().Type == tokenName {
		p, _ := d.GetPhrase()
		if p.Type != phraseNameValue {
			break
		}
		pn := p.Data[0].Data
		pv := p.Data[1].Data
		match := -1
		score := -1
		for i := 0; i < nf; i++ {
			f := t.Field(i)
			if f.PkgPath != "" {
				continue
			}
			tag := f.Tag.Get("ini")
			if tag == "" {
				tag = f.Name
			} else if tag[0] == ',' {
				tag = f.Name + tag
			}
			n, o := parseTag(tag)
			if n == pn {
				match = i
				break
			}
			if l := len(n); l > score && l >= len(pn) && o.Contains("prefix") && pn[:l] == n {
				score = l
				match = i
			}
		}
		if match < 0 {
			continue
		}
		f := s.Field(match)
		switch f.Kind() {
		case reflect.Slice:
			v := reflect.New(f.Type().Elem()).Elem()
			err := setValue(v, pv)
			if err == errUnknownType {
				continue Loop
			} else if err != nil && !d.IgnoreTypeErrors {
				d.Err = err
				return
			}
			reflect.Append(f, v)
		case reflect.Map:
			k := reflect.New(f.Type().Key()).Elem()
			if k.Kind() == reflect.String {
				v := reflect.New(f.Type().Elem()).Elem()
				err := setValue(v, pv)
				if err == errUnknownType {
					continue Loop
				} else if err != nil && !d.IgnoreTypeErrors {
					d.Err = err
					return
				}
				f.SetMapIndex(k, v)
			}
		default:
			v := reflect.New(f.Type()).Elem()
			err := setValue(v, pv)
			if err == errUnknownType {
				continue Loop
			} else if err != nil && !d.IgnoreTypeErrors {
				d.Err = err
				return
			}
			f.Set(v)
		}
	}
}
