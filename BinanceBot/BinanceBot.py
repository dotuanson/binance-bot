import yaml
from binance.client import Client
from binance import ThreadedWebsocketManager
import time
from collections import deque


class BinanceBot(object):
    def __init__(self, symbol):
        with open('./BinanceBot/config.yaml', 'r') as file:
            config = yaml.safe_load(file)
        self.client = Client(config["API"], config["SECRET"])
        self.btc_price = None
        self.bsm = ThreadedWebsocketManager()
        self.start_buy = 0
        self.time_start = time.time()
        self.deq = deque()
        self.symbol = symbol

    def future(self):
        def trade_history(msg):
            ''' define how to process incoming WebSocket messages '''
            error = True if 'error' in msg.values() else False
            if not error:
                text = f"{self.symbol} - {msg['data']['b']} -"
                # print(text, end='\n\r')
                interval = time.time() - self.time_start
                if interval > 1:
                    if self.start_buy > 0:
                        delta = (float(msg['data']['b']) - self.start_buy) / float(msg['data']['b']) * 100
                        self.deq.append(delta)
                        if max(self.deq) <= 0:
                            print(text, "Crash", end='\r')
                        else:
                            print(text, "Bull-up", end='\r')
                        self.deq.clear()
                        time.sleep(0.5)
                    if len(self.deq) == 10:
                        self.deq.popleft()
                    self.start_buy = float(msg['data']['b'])
                    self.time_start = time.time()

        self.bsm.start()
        self.bsm.start_symbol_ticker_futures_socket(callback=trade_history, symbol=self.symbol)

    def spot(self):
        def trade_history(msg):
            text = f"{self.symbol} - {msg['b']}"
            print(text)
        self.bsm.start()
        self.bsm.start_symbol_ticker_socket(callback=trade_history, symbol=self.symbol)
