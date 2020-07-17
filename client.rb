require 'socket'


def send_request
  host = 'localhost'
  port = 9394

  socket = TCPSocket.new host, port
  socket.write("POST /send-metrics HTTP/1.1\r\n")
  socket.write("Transfer-Encoding: chunked\r\n")
  socket.write("Host: #{host}\r\n")
  socket.write("Connection: Close\r\n")
  socket.write("Content-Type: application/octet-stream\r\n")
  socket.write("\r\n")

  message1 = '{"type":"counter","help":"order links","name":"order_links","keys":{"link":"wp_shipping_formatted_address"},"value":1,"prometheus_exporter_action":"increment","custom_labels":{"app":"workers","env":"production"}}'
  message2 = '{"type":"counter","help":"order links","name":"order_links","keys":{"link":"merchant_id"},"value":1,"prometheus_exporter_action":"increment","custom_labels":{"app":"workers","env":"production"}}'

  [message1, message2].each do |message|
    socket.write(message.bytesize.to_s(16).upcase)
    socket.write("\r\n")
    socket.write(message)
    socket.write("\r\n")
    sleep 1
  end

  sleep 1

  socket.write("0\r\n")
  socket.write("\r\n")
  socket.flush
  socket.close
end

trds = []
50.times do
  trds << Thread.new { send_request }
end

trds.map(&:join)