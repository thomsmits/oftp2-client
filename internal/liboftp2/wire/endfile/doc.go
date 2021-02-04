/*
4.5.  End File Phase

4.5.1.  Protocol Sequence

   The Speaker notifies the Listener that it has finished sending a
   Virtual File by sending an End File (EFID) command.  The Listener
   replies with a positive or negative End File command and has the
   option to request a Change Direction command from the Speaker.

   1. Speaker  -- EFID ------------> Listener   End File
               <------------ EFPA --            Answer YES

   2. Speaker  -- EFID ------------> Listener   End File
               <------------ EFPA --            Answer YES + CD
               -- CD -------------->            Change Direction
      Listener <------------ EERP -- Speaker    End to End Response
               -------------- RTR ->            Ready to Receive
      Listener <------------ NERP -- Speaker    Negative End Response
               -------------- RTR ->            Ready to Receive
               Go to Start File Phase

   3. Speaker  -- EFID ------------> Listener   End File
               <------------ EFNA --            Answer NO
*/
package endfile
