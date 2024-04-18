from sys import argv
import random
from hashlib import sha256


P = 0x00c037c37588b4329887e61c2da3324b1ba4b81a63f9748fed2d8a410c2fc21b1232f0d3bfa024276cfd88448197aae486a63bfca7b8bf7754dfb327c7201f6fd17fd7fd74158bd31ce772c9f5f8ab584548a99a759b5a2c0532162b7b6218e8f142bce2c30d7784689a483e095e701618437913a8c39c3dd0d4ca3c500b885fe3

g = 2

if (len(argv) != 3):
    print("usage: python diffy.py [username] [password]")
    exit()

username = argv[1]
password = argv[2]



bigSaltGuy = random.getrandbits(64)

preShaKey = username + ":" + password + ":" + str(bigSaltGuy)


postShaKey = sha256(preShaKey.encode("ascii")).hexdigest()

X = int(postShaKey, 16)

Y = pow(g, X, mod=P)

YAsHex = hex(Y)

outputFileName = input("Please enter a file name: ")

publicKeyFile = open(outputFileName, "w")
publicKeyFile.write(YAsHex)
publicKeyFile.close()

inputPublicKey = input("Please provide public key path: ")

publicKeyFile = open(inputPublicKey, "r")

altKey = publicKeyFile.read()
publicKeyFile.close()

K = pow(Y, int(altKey, 16), mod=P)

outputSymetricKey = sha256(str(K).encode("ascii")).hexdigest()

print(outputSymetricKey)




