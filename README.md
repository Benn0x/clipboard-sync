# Clipboard Sync
### What is it?
It is a smart tool programmed in the popular **Go** language which allows for synchronization of clipboard content between devices on a local network.

### How do you use it?
Simply navigate to [downloads](https://github.com/Benn0x/clipboard-sync/releases/latest) and select the file listed that matches your OS.

Once you've installed it, run the file and it'll prompt you for a passphrase. The passphrase is used to encrypt data so that pesky MITM attacks can be avoided. Please note that the passphrase **must** be the same on both devices. You must repeat these steps for both devices and you will have text appear in the terminal that states what the Host IP is.

### How does it work?
We use a variety of systems such as; MDNs, Websockets, and Encryption. MDNs allows for local device discovery (such as how an Alexa would show in Spotify). Websockets allow for quick and seamless communication between the devices. Encryption means that if you are on a public network people cannot view your clipboard and it's contents.
