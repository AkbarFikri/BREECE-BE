![Alt text](/public/Breece-Banner-Github.png "a title")

Application used to support educational activities within the Brawijaya University environment. We want to help fellow academics inside and outside Brawijaya University to access and take part in various educational activities within the Brawijaya University environment.

<u>**Kelompok 12 BCC Intern**</u>

## 📒 Index

- [📒 Index](#-index)
- [🔰 About](#-about)
- [⚡ Quick Start](#-quick-start)
  - [🔌 Installation](#-installation)
  - [📦 Commands](#-commands)
- [🔧 Development](#-development)
  - [📓 Tech Stack](#-tech-stack)
  - [🔩 API Documentation](#-api-documentation)
  - [📁 File Structure](#-file-structure)
- [🌟 Credit](#-credit)
- [🔒License](#license)

## 🔰 About

Here we want to help fellow academics inside and outside Brawijaya University to access and take part in various educational activities within the Brawijaya University environment. Our application offers various features including:

- **Payment and Registration**, users can easily register and pay activity fees.
- **Email notification**, the user will receive notification of the event they have booked.
- **Upload activities**, users who have the role of organizer will be able to upload their own events.

## ⚡ Quick Start

Here's the step for installation and start our app.

_`Note: This is just a backend Apps not include the Frontend Apps.`_

### 🔌 Installation

1. First, make sure that the go language version you have is more than `1.20`
2. Next, you can clone this repository with the command below

```
$ git clone https://github.com/AkbarFikri/BREECE-BE .
```

3. Provide all the `.env.example` file then rename to `.env`
4. Download all packages needed by Go by running the command below

```
$ go mod tidy
```

**❗ YEAYY Installation Finish!!**

### 📦 Commands

- To run the application you can directly open `main.go` in folder `cmd/app` then click the `run without debugging` button in the right corner of vscode or run the command below

```
$ go run cmd/app/main.go
```

## 🔧 Development

Here is a description of our apps development

### 📓 Tech Stack

List all the Tech Stack we use to build the system in this this project.

| No  | Tech          | Details                                                           |
| --- | ------------- | ----------------------------------------------------------------- |
| 1   | Midtrans      | To provide all payment transaction feature                        |
| 2   | Go            | To build a fast and easy Backend App                              |
| 3   | ID Cloud Host | To provide all application needs related to server infrastructure |
| 4   | Swagger       | To build beatiful documentation                                   |
| 5   | SMTP          | Use for sending all email to user                                 |

### 🔩 API Documentation

- [Swagger](https://breece.akbarfikri.site/api/v1/docs)

_Note : If you have question about the documentation feel free to send message to me._

### 📁 File Structure

Here is our File Structure

```
├───.github
│   └───workflows
├───api
│   └───dist
├───cmd
│   └───app
├───internal
│   ├───app
│   │   ├───config
│   │   ├───entity
│   │   ├───handler
│   │   │   └───rest
│   │   │       ├───middleware
│   │   │       └───routes
│   │   ├───repository
│   │   └───service
│   └───pkg
│       ├───gocron
│       ├───helper
│       ├───mailer
│       │   └───template
│       ├───model
│       └───supabase
└───public
```

## 🌟 Credit

The Member of our team

1. Akbar Fikri Abdillah ( Backend Developer )
2. Pande Gede Natha Satvika ( Product Design )
3. Faidz Gunawan ( Product Manager )
4. Kevin Joshua Silalahi ( Frontend Developer )

## 🔒License

© MIT License - Copyright (c) 2024 AkbarFikri
