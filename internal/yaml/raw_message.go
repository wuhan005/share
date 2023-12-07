// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package yaml

type RawMessage struct {
	unmarshal func(interface{}) error
}

func (msg *RawMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.unmarshal = unmarshal
	return nil
}

func (msg *RawMessage) Unmarshal(v interface{}) error {
	return msg.unmarshal(v)
}
