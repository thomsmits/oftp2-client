/*
5.3.  Command Formats

   The ODETTE-FTP commands are described below using the following
   definitions.

   Position (Pos)

      Field offset within the Command Exchange Buffer, relative to a
      zero origin.

   Field

      The name of the field.

   Description

      A description of the field.

   Format

      F    - A field containing fixed values.  All allowable values for
             the field are enumerated in the command definition.

      V    - A field with variable values within a defined range.  For
             example, the SFIDLRECL field may contain any integer value
             between 00000 and 99999.

      X(n) - An alphanumeric field of length n octets.

        A String contains alphanumeric characters from the following
        set:

         The numerals:               0 to 9
         The upper case letters:     A to Z
         The following special set:  / - . & ( ) space.

        Space is not allowed as an embedded character.

      9(n) - A numeric field of length n octets.

      U(n) - A binary field of length n octets.

             Numbers encoded as binary are always unsigned and in
             network byte order.

      T(n) - An field of length n octets, encoded using [UTF-8].

      String and alphanumeric fields are always left justified and right
      padded with spaces where needed.

      Numeric fields are always right justified and left padded with
      zeros where needed.
*/
package wire
