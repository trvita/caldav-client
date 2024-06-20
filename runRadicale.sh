#! /bin/bash

python3 -m pip install --upgrade radicale
python3 -m radicale --config=./testRadicale/testconfig

#  http://localhost:5232





# https://radicale.org/v3.html

# Radicale is a small but powerful CalDAV (calendars, to-do lists) and CardDAV (contacts) server, that:
# Shares calendars and contact lists through CalDAV, CardDAV and HTTP.
# Supports events, todos, journal entries and business cards.
# Works out-of-the-box, no complicated setup or configuration required.
# Can limit access by authentication.
# Can secure connections with TLS.
# Works with many CalDAV and CardDAV clients.
# Stores all data on the file system in a simple folder structure.
# Can be extended with plugins.
# Is GPLv3-licensed free software.