# This script runs a full end-to-end functional test of the dispatcher and the Shadow transport, using two netcat instances as the application server and application client.
# An alternative way to run this test is to run each command in its own terminal. Each netcat instance can be used to type content which should appear in the other.
FILENAME=testReducedMicroformatsOutput.txt
# Update and build code
go get -u github.com/OperatorFoundation/shapeshifter-dispatcher

# remove text from the output file
rm $FILENAME

# Test -transport
echo "* transport"

~/go/bin/shapeshifter-dispatcher -transparent -server -state state -orport 127.0.0.1:3333 -transport shadow -bindaddr shadow-127.0.0.1:2222 -optionsFile ../../ConfigFiles/ shadowServerChaCha.json -logLevel DEBUG -enableLogging &

sleep 1

killall shapeshifter-dispatcher

#Test -bindhost, -bindport
echo "* bindhost and bindport"
~/go/bin/shapeshifter-dispatcher -transparent -server -state state -orport 127.0.0.1:3333 -transport shadow -bindhost 127.0.0.1 -bindport 2222 -optionsFile ../../ConfigFiles/ shadowServerChaCha.json -logLevel DEBUG -enableLogging &

sleep 1

killall shapeshifter-dispatcher

# Test -targethost, targetport
echo "* targethost and targetport"
~/go/bin/shapeshifter-dispatcher -transparent -client -state state -targethost 127.0.0.1 -targetport 2222 -transport shadow -proxylistenaddr 127.0.0.1:1443 -optionsFile ../../ConfigFiles/ shadowClient.json -logLevel DEBUG -enableLogging &

sleep 1

killall shapeshifter-dispatcher

# Test -proxylistenhost, -proxylistenport
echo "* proxylistenhost and proxylistenport"

~/go/bin/shapeshifter-dispatcher -transparent -client -state state -target 127.0.0.1:2222 -transport shadow -proxylistenhost 127.0.0.1 -proxylistenport 1443 -optionsFile ../../ConfigFiles/ shadowClient.json -logLevel DEBUG -enableLogging &

sleep 1

killall shapeshifter-dispatcher

# Test -mode transparent-TCP
echo "* TransparentTCP"

~/go/bin/shapeshifter-dispatcher -mode transparent-TCP -client -state state -target 127.0.0.1:2222 -transport shadow -proxylistenaddr 127.0.0.1:1443 -optionsFile ../../ConfigFiles/ shadowClient.json -logLevel DEBUG -enableLogging &

sleep 1

killall shapeshifter-dispatcher

# Test -mode transparent-UDP
echo "* TransparentUDP"

~/go/bin/shapeshifter-dispatcher -mode transparent-UDP -client -state state -target 127.0.0.1:2222 -transport shadow -proxylistenaddr 127.0.0.1:1443 -optionsFile ../../ConfigFiles/ shadowClient.json -logLevel DEBUG -enableLogging &

sleep 1

killall shapeshifter-dispatcher

# Test -mode socks5
echo "* socks5"

export TOR_PT_ORPORT=127.0.0.1:2222
~/go/bin/shapeshifter-dispatcher -mode socks5 -client -state state -transport shadow -proxylistenaddr 127.0.0.1:1443 -optionsFile ../../ConfigFiles/ shadowClient.json -logLevel DEBUG -enableLogging &

sleep 1

killall shapeshifter-dispatcher

# Test -mode STUN
echo "* STUN"

~/go/bin/shapeshifter-dispatcher -mode STUN -client -state state -target 127.0.0.1:2222 -transport shadow -proxylistenaddr 127.0.0.1:1443 -optionsFile ../../ConfigFiles/ shadowClient.json -logLevel DEBUG -enableLogging &

sleep 1

killall shapeshifter-dispatcher

OS=$(uname)

if [ "$OS" = "Darwin" ]
then
  FILESIZE=$(stat -f%z "$FILENAME")
else
  FILESIZE=$(stat -c%s "$FILENAME")
fi

if [ "$FILESIZE" = "0" ]
then
  echo "Test Failed"
  killall shapeshifter-dispatcher
  exit 1
fi

echo "Testing complete. Killing processes."

killall shapeshifter-dispatcher
killall nc

echo "Done."
