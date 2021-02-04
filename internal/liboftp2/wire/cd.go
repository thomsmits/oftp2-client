package wire

/*
5.3.12.  CD - Change Direction

   o-------------------------------------------------------------------o
   |       CD          Change Direction                                |
   |                                                                   |
   |       Start File Phase           Speaker ----> Listener           |
   |       End File Phase             Speaker ----> Listener           |
   |       End Session Phase        Initiator <---> Responder          |
   |-------------------------------------------------------------------|
   | Pos | Field     | Description                           | Format  |
   |-----+-----------+---------------------------------------+---------|
   |   0 | CDCMD     | CD Command, 'R'                       | F X(1)  |
   o-------------------------------------------------------------------o

   CDCMD     Command Code                                      Character

      Value: 'R'  CD Command identifier.
*/

// CD requests a change in direction
type CD struct {
}

// CDCMD is the command indicator for the CD command.
const CDCMD = "R"

func (s *CD) FormatDefinition() []FormatDefinition {
	return []FormatDefinition{
		{`F X(1)`, `CDCMD`, &[]Value{{CDCMD, "CD Command"}}},
	}
}

// dataFormat returns the format definition for this command in a machine
// readable form.
func (s *CD) dataFormat() []DataFormat {
	return FormatDefinitionsToDataFormats(s.FormatDefinition())
}

func (s *CD) Command() Command {
	return Command{
		Format: s.dataFormat(),
		Data:   []interface{}{CDCMD},
	}
}

func (s *CD) Marshal() []byte {
	c := s.Command()
	return c.Marshal()
}

func (s *CD) Parse(input []byte) error {
	_, err := NewBuffer(&input, CDCMD)
	if err != nil {
		return err
	}

	return nil
}
