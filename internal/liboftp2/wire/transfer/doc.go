/*
4.4.  Data Transfer Phase

   Virtual File data flows from the Speaker to the Listener during the
   Data Transfer phase, which is entered after the Start File phase.

4.4.1.  Protocol Sequence

   To avoid congestion at the protocol level, a flow control mechanism
   is provided via the Set Credit (CDT) command.

   A Credit limit is negotiated in the Start Session phase; this
   represents the number of Data Exchange Buffers that the Speaker may
   send before it is obliged to wait for a Credit command from the
   Listener.

   The available credit is initially set to the negotiated value by the
   Start File positive answer, which acts as an implicit Credit command.
   The Speaker decreases the available credit count by one for each data
   buffer sent to the Listener.

   When the available credit is exhausted, the Speaker must wait for a
   Credit command from the Listener; otherwise, a protocol error will
   occur and the session will be aborted.

   The Listener should endeavour to send the Credit command without
   delay to prevent the Speaker blocking.

   1. Speaker  -- SFID ------------> Listener   Start File
               <------------ SFPA --            Answer YES

   2. If the credit value is set to 2

      Speaker  -- Data ------------> Listener   Start File
               -- Data ------------>
               <------------- CDT --            Set Credit
               -- Data ------------>
               -- EFID ------------>            End File

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
package transfer
