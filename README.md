# GHODRAT

> WebRTC media servers stress testing tool (currently only Janus)

## COCEPTS

### WebRTC (web Real Time Communication)

> it's a low latency API to exchange video and audio in an efficient manner
> which will enable rich direct peer to peercommunicaton between clients
>
> in the process of connecting 2 peers we need signaling to communicate
> session information with each other (like whatsapp QR, HTTP fetch, WebSocket connection and ...)
> and after signaling is done , 2 peers can connect together via most optimal path

### NAT (Network Address Translation)

> almost every device on internet is behind a NAT, means your are using a router (gateway)
> and your router's IP will communicate with world-wide-web, means first your device
> search for your requested IP address (you google sth and domain-name-server get you acual IP ),
> then your device search that if requested IP is under the same  `subnet` but if not, then it
> will request to the router to transer request packet to requested IP web-server (in the 
> process of transering, the router will create a table which iclude internal nd exernal IP & Port
> which will be replcaed via the router as requester)

<p align="center"><img src="assets/nat.png" /></p>

### NAT Translation Methods

1. One to One NAT (Full-cone NAT)
    - packets to external <IP:Port> on the router always maps to internal <IP:Port> without exception

    <p align="center"><img src="assets/one-to-one-nat.png" /></p>

2. Address restricted NAT
    - packets to external <IP:Port> on the router always maps to internal <IP:Port> as long
    as source address from packet matches the table (regardless of port)

    - allow if we commnicated with this <host> before

    <p align="center"><img src="assets/address-restricted-nat.png" /></p>

3. Port restricted NAT
    - packets to external <IP:Port> on the router always maps to internal <IP:Port> as long
    as source address and port from packet matches the table

    - allow if we commnicated with this <host:port> before


    <p align="center"><img src="assets/port-restricted-nat.png" /></p>

4. symetric NAT
    - packets to external <IP:Port> on the router always maps to internal <IP:Port> as long
    as source address and port from packet matches the table

    - allow if the full pair match

    <p align="center"><img src="assets/symetric-nat.png" /></p>
