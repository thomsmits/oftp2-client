/*
4.2.3.  Secure Authentication

   Having exchanged SSIDs, the Initiator may optionally begin an
   authentication phase, in which each party proves its identity to the
   other.

4.2.4.  Protocol Sequence

   The first authentication message must be sent by the Initiator.

   1. Initiator -- SECD ------------> Responder   Change Direction
                <------------ AUCH --             Challenge
                -- AURP ------------>             Response
                <------------ SECD --             Change Direction
                -- AUCH ------------>             Challenge
                <------------ AURP --             Response

   The Initiator sends a Security Change Direction (SECD) to which the
   Responder replies with an Authentication Challenge (AUCH).

   The Responder looks up the public certificate that is linked to the
   purported identity of the Initiator (located in the SSID).  If the
   Responder is unable to locate a suitable certificate then
   authentication fails.  The Responder uses the public key contained in
   the certificate to encrypt a random challenge, unique for each
   session, for the Initiator.  This encrypted challenge is sent as a
   [CMS] envelope to the Initiator as part of the AUCH.

   The Initiator decrypts the challenge using their private key and
   sends the decrypted challenge back to the Responder in the
   Authentication Response (AURP).

   The Responder checks that the data received in the AURP matches the
   random challenge that was sent to the Initiator.

   If the data matches, then the Initiator has authenticated
   successfully and the Responder replies with a Security Change
   Direction (SECD) beginning the complementary process of verifying the
   Responder to the Initiator.  If the data does not match, then the
   Initiator fails authentication.
*/
package authentication
