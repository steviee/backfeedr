  
**Backfeedr**

Self-Hosted Crash Reporting & App Metrics for iOS Indie Devs

Product Requirements Document

Version 1.1 — March 2026 (Security Update)

Open Source • Single User • Multi-App • Go \+ HTMX \+ SQLite • Privacy-First

# **1\. Vision & Positionierung**

Backfeedr ist ein selbst-gehostetes, leichtgewichtiges Crash-Reporting- und App-Metrik-System für iOS-Indie-Entwickler. Es folgt der Philosophie von Fusionaly: Ein Docker-Container, eine SQLite-Datei, volle Datenkontrolle. Kein Vendor-Lock-in, keine monatlichen Kosten, keine Daten bei Dritten.

## **1.1 Warum Backfeedr?**

| Problem | Bestehende Lösungen | Backfeedr |
| :---- | :---- | :---- |
| Crash Reporting | Sentry (komplex, teuer), Crashlytics (Google-Abhängigkeit), GlitchTip (3+ Container, PostgreSQL) | 1 Container, SQLite, Privacy-first |
| App Metriken | Firebase Analytics (Tracking-Moloch), Countly (Enterprise-Features hinter Paywall) | Nur das Nötigste: Sessions, DAU, Retention |
| Self-Hosting | Sentry Self-Hosted braucht Kafka \+ ClickHouse \+ PostgreSQL \+ Redis | curl | bash, läuft auf 256MB RAM VPS |
| Zielgruppe | Enterprise-fokussiert, überladen für Solo-Devs | Von einem Indie Dev für Indie Devs |

## **1.2 Designprinzipien**

* Single Binary, Single Container — Go-Binary \+ eingebettete Assets, SQLite als einziger Datenspeicher

* Fusionaly-Philosophie — Ein curl-Befehl zum Installieren, Auto-Updates, Backups \= Datei kopieren

* Privacy by Default — Keine persönlichen Daten, keine Cookies, DSGVO-konform ohne Cookie-Banner

* iOS-first — Swift SDK als SPM Package, optimiert für SwiftUI-Apps

* Indie-Dev-Scale — Designed für 1–10 Apps mit jeweils bis zu \~10k DAU

# **2\. Architektur**

## **2.1 Systemübersicht**

`[iOS App] ── HTTPS/JSON ──> [Backfeedr Server (Go)] ──> [SQLite]`

                                    `│`

                              `[HTMX Dashboard]`

                              `[JSON API]`

Der Server ist ein einzelner Go-Binary der sowohl die Ingestion-API als auch das HTMX-Dashboard bedient. SQLite wird über WAL-Mode betrieben für gleichzeitige Lese-/Schreibzugriffe.

## **2.2 Tech Stack**

| Komponente | Technologie | Begründung |
| :---- | :---- | :---- |
| Backend | Go 1.22+ | Single binary, exzellente Performance, geringer RAM |
| HTTP Router | net/http (stdlib) oder chi | Keine unnötigen Dependencies |
| Datenbank | SQLite (WAL-Mode) via modernc.org/sqlite | Pure Go, kein CGO, ein File |
| Dashboard | Go Templates \+ HTMX \+ Alpine.js | Server-rendered, kein Build-Step, kein Node |
| CSS | Pico CSS oder MVP.css | Classless, minimalistisch, Dark Mode |
| iOS SDK | Swift 6, SPM | On-device Symbolication, async/await |
| Container | Docker (Alpine-based) | \~20MB Image |
| Config | Env-Variablen \+ YAML | 12-Factor-App |

## **2.3 Datenbank-Schema (SQLite)**

Alle Tabellen leben in einer einzigen SQLite-Datei. Retention-Policy per Cronjob/Goroutine (Standard: 90 Tage).

### **apps**

| Spalte | Typ | Beschreibung |
| :---- | :---- | :---- |
| id | TEXT (ULID) | Primärschlüssel |
| name | TEXT | App-Name (z.B. "DopaLoop") |
| bundle\_id | TEXT | Bundle Identifier |
| api\_key | TEXT | Ingestion API Key |
| created\_at | DATETIME | Erstellungszeitpunkt |

### **crashes**

| Spalte | Typ | Beschreibung |
| :---- | :---- | :---- |
| id | TEXT (ULID) | Primärschlüssel |
| app\_id | TEXT | FK → apps |
| group\_hash | TEXT | Hash für Crash-Gruppierung |
| exception\_type | TEXT | z.B. EXC\_BAD\_ACCESS |
| exception\_reason | TEXT | Crash-Beschreibung |
| stack\_trace | TEXT (JSON) | Symbolisierter Stack Trace |
| app\_version | TEXT | CFBundleShortVersionString |
| build\_number | TEXT | CFBundleVersion |
| os\_version | TEXT | iOS Version |
| device\_model | TEXT | z.B. iPhone15,2 |
| locale | TEXT | de\_DE |
| free\_memory\_mb | INTEGER | Verfügbarer RAM |
| free\_disk\_mb | INTEGER | Verfügbarer Speicher |
| battery\_level | REAL | 0.0–1.0 |
| is\_charging | BOOLEAN | Ladezustand |
| occurred\_at | DATETIME | Crash-Zeitpunkt |
| received\_at | DATETIME | Server-Empfang |

### **events**

| Spalte | Typ | Beschreibung |
| :---- | :---- | :---- |
| id | TEXT (ULID) | Primärschlüssel |
| app\_id | TEXT | FK → apps |
| type | TEXT | session\_start, session\_end, error, custom |
| name | TEXT | Event-Name (optional) |
| properties | TEXT (JSON) | Beliebige Key-Value-Paare |
| app\_version | TEXT | App-Version |
| os\_version | TEXT | iOS-Version |
| device\_model | TEXT | Gerätemodell |
| session\_id | TEXT | Session-Zuordnung |
| occurred\_at | DATETIME | Zeitpunkt |

### **daily\_metrics (materialisierte Aggregation)**

| Spalte | Typ | Beschreibung |
| :---- | :---- | :---- |
| app\_id | TEXT | FK → apps |
| date | DATE | Tag |
| sessions | INTEGER | Anzahl Sessions |
| unique\_devices | INTEGER | DAU (via anonymem Device-Hash) |
| crashes | INTEGER | Crash-Anzahl |
| errors | INTEGER | Non-Fatal Errors |
| avg\_session\_sec | REAL | Durchschnittliche Session-Dauer |

# **3\. API-Design**

Alle Ingestion-Endpoints akzeptieren JSON und authentifizieren per API-Key im Header. Responses sind minimal, um den Client nicht zu blockieren.

## **3.1 Ingestion API**

| Methode | Endpoint | Beschreibung |
| :---- | :---- | :---- |
| POST | /api/v1/crashes | Crash-Report einliefern |
| POST | /api/v1/events | Event(s) einliefern (Batch) |
| POST | /api/v1/events/batch | Bulk-Event-Upload (bis zu 100 Events) |
| GET | /api/v1/health | Health-Check (kein Auth) |

## **3.2 Authentifizierung**

`X-Backfeedr-Key: bf_live_a1b2c3d4e5f6...`

Ein API-Key pro App. Keys haben das Präfix bf\_live\_ für Produktion und bf\_test\_ für Entwicklung. Test-Keys werden im Dashboard separat dargestellt.

## **3.3 Crash-Report Payload**

`{`

  `"exception_type": "EXC_BAD_ACCESS",`

  `"exception_reason": "Attempted to dereference null pointer",`

  `"stack_trace": [`

    `{ "frame": 0, "symbol": "ContentView.body.getter",`

      `"file": "ContentView.swift", "line": 42 },`

    `{ "frame": 1, "symbol": "SwiftUI.View.update()",`

      `"file": null, "line": null }`

  `],`

  `"app_version": "1.2.0",`

  `"build_number": "47",`

  `"os_version": "18.3.1",`

  `"device_model": "iPhone16,1",`

  `"device_id_hash": "a9f3...c721",`

  `"locale": "de_DE",`

  `"free_memory_mb": 312,`

  `"battery_level": 0.67,`

  `"occurred_at": "2026-03-12T14:22:31Z"`

`}`

## **3.4 Event Payload**

`{`

  `"type": "session_start",`

  `"session_id": "sess_a1b2c3",`

  `"app_version": "1.2.0",`

  `"os_version": "18.3.1",`

  `"device_model": "iPhone16,1",`

  `"device_id_hash": "a9f3...c721",`

  `"properties": { "source": "widget" },`

  `"occurred_at": "2026-03-12T14:22:31Z"`

`}`

# **4\. Swift SDK (BackfeedrKit)**

## **4.1 Design-Ziele**

* SPM-only Distribution (kein CocoaPods, kein Carthage)

* Swift 6 Concurrency mit Sendable-Conformance

* Einzeiler-Setup im App-Einstiegspunkt

* On-device Symbolication via SPI oder PLCrashReporter

* Offline-Queue: Crashes werden lokal gespeichert und beim nächsten Start gesendet

* Zero persönliche Daten: Device-ID wird als irreversibler SHA-256-Hash gespeichert

* Minimaler Footprint: Kein Impact auf App-Startup-Zeit

## **4.2 Integration**

Package.swift:

`.package(url: "https://github.com/backfeedr/backfeedr-swift", from: "1.0.0")`

App-Setup (SwiftUI):

`import BackfeedrKit`

`@main`

`struct DopaLoopApp: App {`

    `init() {`

        `Backfeedr.configure(`

            `endpoint: "https://crashes.meinserver.de",`

            `apiKey: "bf_live_a1b2c3d4"`

        `)`

    `}`

    `var body: some Scene {`

        `WindowGroup { ContentView() }`

    `}`

`}`

## **4.3 API-Oberfläche**

`// Automatisch: Crash-Handling, Session-Tracking`

`// Manuell: Non-Fatal Errors`

`Backfeedr.capture(error: myError)`

`Backfeedr.capture(error: myError, context: ["screen": "settings"])`

`// Custom Events`

`Backfeedr.track("purchase_completed", properties: ["plan": "pro"])`

`// Breadcrumbs (letzte 20 vor Crash)`

`Backfeedr.breadcrumb("Tapped settings button")`

## **4.4 Datenschutz-Architektur**

Das SDK sammelt absichtlich keine persönlich identifizierbaren Daten:

| Datenpunkt | Methode | Zweck |
| :---- | :---- | :---- |
| Device-ID | SHA-256(identifierForVendor) | DAU-Zählung ohne Tracking |
| Gerätemodell | sysctl hw.machine | Crash-Kontext |
| iOS-Version | UIDevice.current.systemVersion | Kompatibilitäts-Analyse |
| App-Version | CFBundleShortVersionString | Release-Tracking |
| Locale | Locale.current.identifier | Regionale Crash-Muster |
| RAM/Disk | ProcessInfo / FileManager | OOM-Korrelation |

Was das SDK NICHT sammelt: IP-Adressen (der Server loggt keine IPs), Nutzernamen, E-Mail-Adressen, Standortdaten, Bildschirminhalte, Tastatureingaben.

# **5\. Dashboard (HTMX)**

## **5.1 Views**

| View | Inhalt | Prio |
| :---- | :---- | :---- |
| Übersicht | Crash-Free Rate, DAU, Sessions (7d/30d Sparklines), Top-Crashes | MVP |
| Crashes | Gruppierte Crash-Liste, Occurrence-Count, Affected Versions, Stack Trace Detail | MVP |
| Crash Detail | Vollständiger Stack Trace, Device-Verteilung, Timeline, Breadcrumbs | MVP |
| Events | Event-Stream mit Filtern (Typ, App, Version, Zeitraum) | MVP |
| Metriken | DAU/MAU, Session-Dauer, Retention (D1/D7/D30), Version-Adoption | v1.1 |
| Apps | App-Verwaltung, API-Keys generieren/rotieren | MVP |
| Einstellungen | Retention-Policy, Auth-Token, Export | MVP |

## **5.2 Dashboard-Technologie**

* Go html/template mit eingebetteten Templates (embed.FS)

* HTMX für partielle Updates (Crash-Liste, Filter, Pagination)

* Alpine.js für leichte Interaktivität (Dropdowns, Toggles)

* Chart.js (via CDN) für Sparklines und Trends

* Pico CSS für classless Styling mit Dark Mode

* Kein Build-Step, kein npm, kein Webpack

## **5.3 Authentifizierung**

Single-User-System: Ein Bearer-Token wird beim ersten Start generiert und in der Konfigurationsdatei gespeichert. Das Dashboard ist per Token geschützt (Cookie-basiert nach Login). Optional: Basic Auth hinter Reverse Proxy (Caddy/Nginx).

# **6\. Deployment & Operations**

## **6.1 One-Line Install**

`curl -fsSL https://backfeedr.dev/install | bash`

Das Install-Script erledigt: Docker-Image pull, SQLite-Volume erstellen, Caddy/Traefik-Config für HTTPS, systemd-Service anlegen, initiales Admin-Token generieren.

## **6.2 Docker Compose**

`services:`

  `backfeedr:`

    `image: ghcr.io/backfeedr/backfeedr:latest`

    `ports:`

      `- "8080:8080"`

    `volumes:`

      `- ./data:/data`

    `environment:`

      `- BACKFEEDR_AUTH_TOKEN=bf_admin_xxxxx`

      `- BACKFEEDR_BASE_URL=https://crashes.example.com`

      `- BACKFEEDR_RETENTION_DAYS=90`

## **6.3 Konfiguration**

| Variable | Default | Beschreibung |
| :---- | :---- | :---- |
| BACKFEEDR\_AUTH\_TOKEN | (generiert) | Admin-Token fürs Dashboard |
| BACKFEEDR\_BASE\_URL | http://localhost:8080 | Externe URL |
| BACKFEEDR\_RETENTION\_DAYS | 90 | Automatische Löschung alter Daten |
| BACKFEEDR\_DB\_PATH | /data/backfeedr.db | SQLite-Datenbankpfad |
| BACKFEEDR\_MAX\_BODY\_SIZE | 1MB | Max. Request-Größe |
| BACKFEEDR\_RATE\_LIMIT | 100/min | Rate-Limit pro API-Key |

## **6.4 Backup & Recovery**

Backup ist trivial — genau wie bei Fusionaly:

`cp /data/backfeedr.db /backup/backfeedr-$(date +%Y%m%d).db`

Für konsistente Backups bei laufendem Server: SQLite .backup-Befehl oder VACUUM INTO.

## **6.5 Systemanforderungen**

|  | Minimum | Empfohlen |
| :---- | :---- | :---- |
| RAM | 256 MB | 512 MB |
| CPU | 1 vCPU (arm64 oder x86) | 1 vCPU |
| Disk | 1 GB | 5 GB (abhängig von Retention) |
| OS | Jedes OS mit Docker | Debian 12 / Ubuntu 24 |

# **7\. Security & Privacy**

Backfeedrs Sicherheitsmodell folgt dem Prinzip der konzentrischen Verteidigungsringe. Die wichtigste Schutzmaßnahme ist gleichzeitig die einfachste: Keine persönlichen Daten sammeln. Wenn ein abgefangener Crash-Report keine PII enthält, ist er für einen Angreifer wertlos.

## **7.1 Schicht 0 — Datenminimierung (Privacy by Design)**

Die beste Verteidigung gegen Datenmissbrauch ist, keine missbrauchbaren Daten zu haben. Backfeedr sammelt ausschließlich technische Daten und implementiert aktive Schutzmaßnahmen gegen versehentliche PII-Leaks:

### **SDK-seitige PII-Scrubbing-Pipeline**

* Automatische Regex-Filterung in Breadcrumbs und Custom Properties vor dem Versand

* Erkannte Muster: E-Mail-Adressen, Telefonnummern, URLs mit Query-Parametern, IPv4/IPv6-Adressen

* Ersetzung durch \[REDACTED\] — geschieht on-device, bevor Daten das Gerät verlassen

* Stack Traces enthalten nur Symbolnamen, Dateinamen und Zeilennummern — niemals Variablenwerte

* Device-ID ist ein irreversibler SHA-256-Hash von identifierForVendor — kein Tracking möglich

### **Server-seitige Maßnahmen**

* IP-Adressen werden weder geloggt noch gespeichert — der Request-Handler verwirft die Source-IP sofort

* Keine Cookies, kein Session-Tracking, kein Fingerprinting auf der Ingestion-API

* Automatische Retention-Policy: Daten werden nach konfigurierbarer Frist (Standard: 90 Tage) gelöscht

* SQLite-Datenbank ist verschlüsselt speicherbar (SQLCipher als optionales Feature)

| Datenpunkt | Gesammelt? | Methode |
| :---- | :---- | :---- |
| Device-Modell | Ja | sysctl hw.machine → z.B. iPhone16,1 |
| iOS-Version | Ja | UIDevice.current.systemVersion |
| App-Version \+ Build | Ja | CFBundleShortVersionString \+ CFBundleVersion |
| Locale | Ja | Locale.current.identifier (z.B. de\_DE) |
| RAM / Disk frei | Ja | ProcessInfo / FileManager |
| Batterielevel | Ja | UIDevice.current.batteryLevel |
| Device-ID | Pseudonymisiert | SHA-256(identifierForVendor) |
| IP-Adresse | Nein | Wird serverseitig sofort verworfen |
| Nutzername / E-Mail | Nein | Wird nicht erhoben |
| Standort | Nein | Kein Zugriff auf Location Services |
| Bildschirminhalte | Nein | Kein Screenshot-Capture |
| Tastatureingaben | Nein | Kein Keylogging |

## **7.2 Schicht 1 — Transport-Sicherheit**

Alle Kommunikation zwischen SDK und Server läuft über TLS 1.3. Das SDK verweigert HTTP-Verbindungen komplett — keine NSAppTransportSecurity-Ausnahmen, kein Fallback auf ältere TLS-Versionen.

* TLS 1.3 als Minimum — kein TLS 1.2 Fallback

* HSTS-Header auf dem Server (Strict-Transport-Security: max-age=63072000; includeSubDomains)

* SDK erzwingt HTTPS — configure() schlägt bei http://-Endpoints fehl

* Empfehlung: Caddy als Reverse Proxy (automatisches HTTPS via Let’s Encrypt)

## **7.3 Schicht 2 — Request-Authentifizierung (HMAC Signing)**

Der API-Key allein identifiziert nur die App — er schützt nicht gegen Replay-Angriffe oder Manipulation. Deshalb signiert das SDK jeden Request mit HMAC-SHA256:

`// SDK-seitig (Swift)`

`let timestamp = ISO8601DateFormatter().string(from: Date())`

`let payload = "\(timestamp).\(requestBodyHash)"`

`let signature = HMAC<SHA256>.authenticationCode(`

    `for: Data(payload.utf8),`

    `using: SymmetricKey(data: Data(apiKey.utf8))`

`)`

`// HTTP Headers`

`X-Backfeedr-Key: bf_live_a1b2c3d4`

`X-Backfeedr-Timestamp: 2026-03-12T14:22:31Z`

`X-Backfeedr-Signature: sha256=a9f3c721...`

Der Server verifiziert: (1) Signatur stimmt, (2) Timestamp liegt innerhalb eines 5-Minuten-Fensters (verhindert Replay-Angriffe), (3) API-Key ist gültig und aktiv.

## **7.4 Schicht 3 — Public Key Pinning (Optional)**

Certificate Pinning verhindert DNS-Hijacking und Proxy-basierte MitM-Angriffe, indem das SDK nur Verbindungen zu Servern akzeptiert, deren Public Key einem bekannten Hash entspricht. Backfeedr pinnt den SPKI-Hash (Subject Public Key Info), nicht das gesamte Zertifikat — so überlebt der Pin auch Zertifikatsrotationen bei Let’s Encrypt.

`// SDK-Konfiguration mit optionalem Pin`

`Backfeedr.configure(`

    `endpoint: "https://crashes.example.com",`

    `apiKey: "bf_live_a1b2c3d4",`

    `publicKeyHashes: [`

        `"sha256/AAAA+BBBB...",  // Primärer Pin`

        `"sha256/CCCC+DDDD..."   // Backup-Pin (nächstes Zertifikat)`

    `]`

`)`

Implementierung via URLSessionDelegate mit URLAuthenticationChallenge. Der Backup-Pin ermöglicht nahtlose Zertifikatsrotation. Dokumentation erklärt, wie der SPKI-Hash per OpenSSL extrahiert wird:

`openssl s_client -connect crashes.example.com:443 | \`

  `openssl x509 -pubkey -noout | \`

  `openssl pkey -pubin -outform der | \`

  `openssl dgst -sha256 -binary | base64`

## **7.5 Schicht 4 — DNS-TXT Endpoint-Validierung (Leichtgewichtige Alternative)**

Als leichtgewichtige Alternative zu Certificate Pinning kann der Server seinen öffentlichen Schlüssel in einem DNS-TXT-Record veröffentlichen. Das SDK verifiziert beim ersten Start (und periodisch danach), dass der Endpoint legitim ist:

`// DNS TXT Record`

`_backfeedr.crashes.example.com TXT "v=bf1; k=sha256/AAAA+BBBB..."`

`// SDK prüft beim configure():`

`// 1. DNS-TXT-Record für _backfeedr.{endpoint-host} abfragen`

`// 2. Public Key Hash aus Record extrahieren`

`// 3. Gegen Server-Zertifikat validieren`

`// 4. Ergebnis cachen (TTL: 24h)`

Vorteil gegenüber hartkodiertem Pinning: Kein App-Update nötig bei Schlüsselrotation. Der DNS-Record wird vom Server-Betreiber aktualisiert, das SDK prüft periodisch. Nachteil: DNS-Queries können selbst manipuliert werden (Mitigation: DNSSEC-Empfehlung in der Dokumentation).

## **7.6 Schicht 5 — Rate-Limiting & Abuse Prevention**

| Maßnahme | Konfiguration | Zweck |
| :---- | :---- | :---- |
| Rate-Limit | 100 Requests/min pro API-Key | Verhindert Flooding/DDoS |
| Request-Size-Limit | Max. 1 MB pro Request | Verhindert Speicher-Exhaustion |
| Timestamp-Window | Max. 5 Min. Abweichung | Verhindert Replay-Angriffe |
| Stale-Report-Rejection | Reports älter als 24h verwerfen | Verhindert historische Injection |
| API-Key-Rotation | Jederzeit im Dashboard | Kompromittierte Keys invalidieren |
| Batch-Limit | Max. 100 Events pro Batch-Request | Verhindert Bulk-Abuse |

## **7.7 Schicht 6 — Apple App Attest (Future)**

Apple’s App Attest Service (Teil des DeviceCheck-Frameworks, seit iOS 14\) ermöglicht kryptografische Verifizierung, dass ein Request von einer legitimen App-Instanz auf einem echten Apple-Gerät stammt. Der Schlüssel wird im Secure Enclave generiert und von Apple verifiziert.

Workflow:

1. 1\. SDK generiert Schlüsselpaar via DCAppAttestService.generateKey() — Private Key bleibt im Secure Enclave

2. 2\. Server sendet Challenge (Nonce) an das SDK

3. 3\. SDK attestiert den Schlüssel via Apple’s Server → Attestation Object

4. 4\. Server validiert Attestation (CBOR-Parsing, X.509-Zertifikatskette gegen Apple Root CA)

5. 5\. Folge-Requests werden mit Assertions signiert und serverseitig verifiziert

Einschränkungen: Die Server-seitige Validierung erfordert CBOR-Decoding und X.509-Zertifikatskettenprüfung in Go — nicht-trivialer Implementierungsaufwand. Apple empfiehlt App Attest als optionale Maßnahme, da Requests an Apple’s Server fehlschlagen können. Daher geplant für v1.2 als Opt-in-Feature.

## **7.8 Sicherheitsarchitektur — Übersicht**

| Schicht | Maßnahme | Aufwand | Phase | Schützt gegen |
| :---- | :---- | :---- | :---- | :---- |
| 0 | PII-Scrubbing, keine IP-Logs | Niedrig | MVP | Daten-Leak bei Kompromittierung |
| 1 | TLS 1.3 only, HSTS | Niedrig | MVP | Passive Abhörangriffe |
| 2 | HMAC Request-Signing | Mittel | MVP | Replay-Angriffe, Request-Manipulation |
| 3 | Public Key Pinning (opt.) | Mittel | v1.0 | DNS-Hijacking, Proxy-MitM |
| 4 | DNS-TXT Validierung (opt.) | Mittel | v1.0 | Endpoint-Spoofing ohne App-Update |
| 5 | Rate-Limiting \+ Timestamps | Niedrig | MVP | Flooding, DDoS, Replay |
| 6 | Apple App Attest (opt.) | Hoch | v1.2 | Fake Clients, App-Clones |

## **7.9 Threat Model**

| Bedrohung | Risiko | Mitigation |
| :---- | :---- | :---- |
| DNS-Redirect auf fremdes Backend | Mittel | Public Key Pinning oder DNS-TXT-Validierung; HMAC-Signatur schlägt auf fremdem Server fehl |
| API-Key aus App extrahiert | Hoch | HMAC-Signing verhindert willkürliche Requests; App Attest verifiziert App-Integrität; Key-Rotation im Dashboard |
| MitM-Proxy fängt Crash-Reports ab | Mittel | TLS 1.3 \+ Public Key Pinning; PII-Scrubbing macht abgefangene Daten wertlos |
| Replay abgefangener Requests | Niedrig | Timestamp-Window (5 Min.) \+ HMAC-Signatur |
| DDoS auf Ingestion-API | Mittel | Rate-Limiting \+ Request-Size-Limit; hinter Caddy/Cloudflare |
| Server-Kompromittierung | Niedrig | Keine PII in DB; verschlüsselte SQLite (optional); regelmäßige Backups |
| Fake Crash-Reports einschleusen | Niedrig | HMAC-Signatur \+ API-Key; App Attest (v1.2) |

# **8\. Roadmap**

| Phase | Features | Zeitrahmen |
| :---- | :---- | :---- |
| MVP (v0.1) | Go-Server mit Ingestion API, SQLite-Storage, Crash-Gruppierung, minimales HTMX-Dashboard (Crash-Liste \+ Detail), Swift SDK mit Crash-Handling \+ Session-Tracking \+ PII-Scrubbing \+ HMAC-Signing, TLS 1.3 only, Rate-Limiting, Docker-Image | 4–6 Wochen |
| v0.2 | Non-Fatal Error Logging, Breadcrumbs, Event-Stream View, API-Key-Rotation, Retention-Cronjob | 2–3 Wochen |
| v1.0 | DAU/MAU Metriken, Retention-Charts (D1/D7/D30), Version-Adoption, Crash-Free Rate, Public Key Pinning (opt.), DNS-TXT-Validierung (opt.), One-Line Installer, README \+ Docs | 3–4 Wochen |
| v1.1 | Alerting (Webhook/Ntfy/Gotify bei neuem Crash-Typ), JSON-Export, Performance-Metriken (App Start Time), Apple Watch SDK, SQLCipher-Verschlüsselung (opt.) | TBD |
| v1.2 | Apple App Attest (opt.), Multi-User (opt.), dSYM-Upload als Alternative zu On-Device Symbolication, Grafana-Integration (optional) | TBD |

# **9\. Projektstruktur**

`backfeedr/`

`├── cmd/`

`│   └── backfeedr/`

`│       └── main.go              # Einstiegspunkt`

`├── internal/`

`│   ├── server/`

`│   │   ├── server.go            # HTTP-Server Setup`

`│   │   ├── middleware.go        # Auth, Rate-Limit, CORS`

`│   │   └── routes.go            # Route-Registration`

`│   ├── api/`

`│   │   ├── crashes.go           # POST /api/v1/crashes`

`│   │   ├── events.go            # POST /api/v1/events`

`│   │   └── health.go            # GET /api/v1/health`

`│   ├── dashboard/`

`│   │   ├── handler.go           # HTMX Page Handlers`

`│   │   ├── templates/           # Go Templates (embed.FS)`

`│   │   └── static/              # CSS, JS, Icons`

`│   ├── store/`

`│   │   ├── sqlite.go            # SQLite-Verbindung + Migration`

`│   │   ├── crashes.go           # Crash CRUD + Gruppierung`

`│   │   ├── events.go            # Event CRUD`

`│   │   └── metrics.go           # Aggregation + daily_metrics`

`│   └── config/`

`│       └── config.go            # Env-Parsing`

`├── Dockerfile`

`├── docker-compose.yml`

`├── go.mod`

`├── go.sum`

`├── CLAUDE.md                        # Claude Code Brief`

`├── LICENSE                          # MIT`

`└── README.md`

Das Swift SDK lebt in einem separaten Repository: backfeedr/backfeedr-swift

# **10\. Crash-Gruppierung**

Crashes werden anhand eines group\_hash gruppiert. Der Hash wird aus den folgenden Feldern berechnet:

* exception\_type (z.B. EXC\_BAD\_ACCESS)

* Top 3 App-Frames im Stack Trace (Symbole, die zum Bundle gehören)

* Nicht einbezogen: System-Frames, Adressen, Zeilennummern

Algorithmus: SHA-256("{exception\_type}:{frame1\_symbol}:{frame2\_symbol}:{frame3\_symbol}"). So werden identische Crashes zusammengefasst, auch wenn sie in leicht unterschiedlichen Umgebungen auftreten.

# **11\. Vergleich mit Alternativen**

|  | Backfeedr | Sentry (Self-Host) | GlitchTip | Crashlytics |
| :---- | :---- | :---- | :---- | :---- |
| Container | 1 | 10+ | 3+ | SaaS only |
| Datenbank | SQLite | PostgreSQL \+ ClickHouse | PostgreSQL | Google Cloud |
| Min. RAM | 256 MB | 4+ GB | 1 GB | — |
| iOS SDK | Eigen (SPM) | Sentry SDK | Sentry SDK | Firebase SDK |
| Privacy | Self-hosted, no PII | Self-hosted | Self-hosted | Google |
| Kosten | Kostenlos (OSS) | Kostenlos (OSS) | Kostenlos (OSS) | Kostenlos (Lock-in) |
| Setup | curl | bash | Komplex | Mittel | SDK \+ Console |
| Zielgruppe | Indie Devs | Enterprise | Small Teams | Alle |

# **12\. Offene Fragen**

* Symbolication-Qualität: Reicht On-Device Symbolication für Release-Builds? Brauchen wir früh eine dSYM-Upload-Option als Fallback?

* Alerting MVP: Soll v1.0 bereits Webhook-Alerts bei neuen Crash-Typen enthalten, oder ist das v1.1?

* Naming: backfeedr.dev als Domain? GitHub-Org backfeedr verfügbar?

* Lizenz: MIT für Server und SDK (indie-freundlich, kein AGPL-Risiko für App-Store-Apps).

* Monetarisierung: Rein OSS, oder Pro-Tier mit z.B. AI-powered Crash Analysis (à la Fusionaly)?

* App Attest Go-Library: Existiert eine ausgereifte Go-Bibliothek für CBOR/X.509-Validierung, oder muss das selbst implementiert werden?

* DNS-TXT vs. Certificate Pinning: Beide als Optionen anbieten, oder nur eines davon im MVP?

* SQLCipher: Lohnt sich Datenbank-Verschlüsselung bei reinen technischen Daten ohne PII, oder ist das Security-Theater?

