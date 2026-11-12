What is the main goal of the application?
To deliver fast, secure clipboard content between multiple devices within the same network.

How would it communicate?
MDNs for device discovery
Websockets for connections to the host device

How would we secure it?
Using RSA we could encrypt data that goes across the network. This would make it near-impossible for malicious users to know what is in another persons clipboard.



MDNs Feature
On MDNs device registry we will search for local devices running the software.
If there is a device we will have the text include host:false
If there is not a device we will have the text say host:true instead.

If the device is a host we should open a websocket connection.

Websocket Feature
On startup of the websocket we should generate a 