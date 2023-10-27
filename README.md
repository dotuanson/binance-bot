## 📋 Table of Contents

* [About the Project](#about-the-project)
  * [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
* [Usage](#usage)
* [Contributing to project](#contributing-to-project)
* [Contributors](#contributors)

<h2 id="about-the-project"> 📑 About the project </h2>

**Binance bot** is a repo building a trading bot using Binance API

<h3 id="built-with"> 💻 Built with </h3>

* [Golang](https://go.dev)

<h2 id="getting-started"> 🛠 Getting started </h2>

<h3 id="prerequisites"> 📎 Prerequisites </h3>

- Go `1.20.2`
- GNU Make `3.81`

<h2 id="installation"> 📎 Installation </h2>

1. Clone the repo
```sh
git clone https://github.com/dotuanson/binance-bot
```
2. Enter your configurations in `.env`
```
COIN_LIST=""
API_KEY=""
SECRET_KEY=""
BASE_URL=""
TELEGRAM_TOKEN=""
TELEGRAM_CHATID=""
```

<h2 id="usage">  🤖 Usage </h2>

### 📌 Golang testing
```commandline
make test
```

### 📌 Golang building
```commandline
make build
```

### 📌 Docker compose
```commandline
make deploy
```



<h2 id="contributing-to-project"> 👋 Contributing to project </h2>

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/<YOUR-FEATURE>`)
3. Commit your Changes (`git commit -m 'Add some <YOUR-FEATURE>'`)
4. Push to the Branch (`git push origin feature/<YOUR-FEATURE>`)
5. Open a Pull Request

<h2 id="contributors"> 👨‍💻 Contributors </h2>
