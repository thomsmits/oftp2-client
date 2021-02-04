/*
4.2.1.  Entity Definition

   The ODETTE-FTP entity that took the initiative to establish the
   network connection becomes the Initiator.  Its peer becomes the
   Responder.

4.2.2.  Protocol Sequence

   The first message must be sent by the Responder.

   1. Initiator <-------------SSRM -- Responder   Ready Message
                -- SSID ------------>             Identification
                <------------ SSID --             Identification

*/
package session
