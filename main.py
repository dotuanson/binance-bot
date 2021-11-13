from BinanceBot.BinanceBot import BinanceBot
import yaml

with open('config_bot.yaml', 'r') as file:
    config = yaml.safe_load(file)


def future():
    bots = list()
    for coin in config["Future"]:
        bots.append(BinanceBot(symbol=coin))
    for bot in bots:
        bot.future()


def spot():
    bots = list()
    for coin in config["Spot"]:
        bots.append(BinanceBot(symbol=coin))
    for bot in bots:
        bot.spot()


if __name__ == '__main__':
    future()
