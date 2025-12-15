# Ambassador System — Complete Specification

> **Version:** 1.0  
> **Last Updated:** December 8, 2025  
> **Status:** Implementation Ready

---

## 🖼️ REQUIRED ASSETS — Image Checklist

**Place all images in:** `services/discordbot/assets/`

| # | Filename (Rename To) | Purpose | Dimensions |
|---|----------------------|---------|------------|
| 1 | `waitlist.png` | Header banner for `/wait-list` embed | 1200×400 |
| 2 | `ourvision.png` | Header banner for `/owls-vision` embed | 1200×400 |
| 3 | `leaderboard.png` | Header banner for `/leaderboard` embed | 1200×400 |
| 4 | `panel.png` | Header banner for `/panel` dashboard | 1200×400 |
| 5 | `ambassadorwelcome.png` | Header banner for ambassador private channel welcome | 1200×400 |
| 6 | `approved.png` | Header banner shown when applicant is approved/accepted | 1200×400 |
| 7 | `conversion.png` | Header banner for conversion success notification | 1200×400 |
| 8 | `footer.png` | Footer image for ALL embeds (can reuse existing `owlsfooter.png`) | 1200×100 |

**Total Images Needed: 8** (or 7 if reusing existing footer)

---

## 🗺️ VISUAL FLOWCHART

> **See companion document:** [`AMBASSADOR_FLOWCHART.md`](./AMBASSADOR_FLOWCHART.md)
> 
> Contains extremely detailed visual flowcharts for:
> - Complete system architecture tree
> - Database schema tree
> - Applicant journey (waitlist → ambassador)
> - `/onboard` command step-by-step
> - Interview accept/reject flows
> - `/claim` validation gauntlet (6 checks)
> - `/unclaim` with cooldown logic
> - Conversion detection engine
> - Payout request & completion flows
> - Leaderboard weekly reset mechanism
> - Panel button interactions
> - Admin balance commands
> - All Custom IDs reference

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Server Configuration](#2-server-configuration)
3. [Role Configuration](#3-role-configuration)
4. [Channel Configuration](#4-channel-configuration)
5. [Database Schema](#5-database-schema)
6. [Asset Requirements](#6-asset-requirements)
7. [Commands Reference](#7-commands-reference)
8. [Embed Specifications](#8-embed-specifications)
9. [Workflows](#9-workflows)
10. [Business Rules](#10-business-rules)
11. [Logging System](#11-logging-system)
12. [Permission Matrix](#12-permission-matrix)
13. [Pre-Implementation Checklist](#13-pre-implementation-checklist)
14. [Technical Implementation Notes](#14-technical-implementation-notes)

---

## 1. System Overview

### 1.1 Purpose

The Ambassador System is a cross-server referral and conversion tracking system designed to:

- Recruit and onboard ambassadors through a structured application process
- Enable ambassadors to "claim" free members in the Main Server
- Track conversions when claimed members upgrade to paid membership
- Manage ambassador commissions and payouts
- Provide leaderboards and performance analytics

### 1.2 Core Components

| Component | Description |
|-----------|-------------|
| **Waitlist System** | Public-facing application entry point via Tally form |
| **Onboarding Pipeline** | Staff-managed interview and approval process |
| **Claims Engine** | Cross-server claiming mechanism with limits and cooldowns |
| **Conversion Tracker** | Automatic detection of claimed user upgrades |
| **Payout System** | Commission balance tracking and payout requests |
| **Leaderboard** | Weekly and lifetime performance rankings |
| **Logging** | Admin audit trail for all ambassador activities |

### 1.3 Multi-Server Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         MAIN SERVER                                  │
│                    (1431718856619200686)                            │
│                                                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────┐  │
│  │ Free Users  │───▶│   /claim    │───▶│ Conversion Detection    │  │
│  │ (claimable) │    │ (amb only)  │    │ (role change monitor)   │  │
│  └─────────────┘    └─────────────┘    └─────────────────────────┘  │
│                                                                      │
└──────────────────────────────┬──────────────────────────────────────┘
                               │
                    Bot operates in BOTH servers
                               │
┌──────────────────────────────▼──────────────────────────────────────┐
│                       AMBASSADOR SERVER                              │
│                    (1441909666308423874)                            │
│                                                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────┐  │
│  │  Waitlist   │───▶│  Onboarding │───▶│  Ambassador Dashboard   │  │
│  │  /wait-list │    │  /onboard   │    │  /panel, /leaderboard   │  │
│  └─────────────┘    └─────────────┘    └─────────────────────────┘  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. Server Configuration

### 2.1 Server Identifiers

| Server | Guild ID | Purpose |
|--------|----------|---------|
| **Main Server** | `1431718856619200686` | Public community, free/paid members, where `/claim` is executed |
| **Ambassador Server** | `1441909666308423874` | Private ambassador hub, onboarding, panels, payouts |

### 2.2 Bot Requirements

The Owls Capital bot **MUST** be present and operational in **BOTH** servers simultaneously.

**Critical:** The bot role must be positioned **at the top of the role hierarchy** in both servers to ensure it can:
- Assign/remove roles to any member
- Create and delete channels/threads
- Manage permissions on channels
- Read member information cross-server

---

## 3. Role Configuration

### 3.1 Ambassador Server Roles (1441909666308423874)

| Role Name | Role ID | Purpose | Position |
|-----------|---------|---------|----------|
| **Staff/Admin** | `1441909798852366499` | Full access to admin commands, can onboard/manage ambassadors | HIGH (near top) |
| **Approved Owl** | `1445249241743364236` | Granted after `/onboard`, pending interview | MEDIUM |
| **Full Ambassador** | `1445251690365456466` | Granted after interview acceptance, unlocks full server access | MEDIUM |
| **Ambassador (Main)** | _TO BE CREATED_ | Role in Main Server for cross-server identification | MEDIUM |

### 3.2 Main Server Roles (1431718856619200686)

| Role Name | Role ID | Purpose | Claimable? |
|-----------|---------|---------|------------|
| **Default Member** | `718643370301325404` | Free community members | ✅ YES |
| **No Role** | — | Users without any role | ✅ YES |
| **Paying Member** | `885910828086362132` | Active paid subscribers | ❌ NO |
| **Lifetime Member** | `718643316786462772` | Lifetime paid access | ❌ NO |
| **Ambassador** | _TO BE CREATED_ | Identifies ambassadors in Main Server | ❌ NO |

### 3.3 Roles to Create

You must manually create these roles before running the system:

#### In Main Server (1431718856619200686):
```
Role Name: Ambassador
Color: (Your choice - recommend gold/amber #FFB800)
Permissions: None special needed
Position: Above Default Member, below Paying Member
Mentionable: No
Hoisted: Optional
```

#### In Ambassador Server (1441909666308423874):
> All required roles appear to be created based on the IDs provided.

---

## 4. Channel Configuration

### 4.1 Ambassador Server Channels (1441909666308423874)

| Channel/Category | ID | Purpose |
|-----------------|-----|---------|
| **Approved Applicants Channel** | `1445245329778806834` | Where approved applicants can see their application status; interview threads created here |
| **Ambassador Private Channels Category** | `1445215664515055779` | Parent category for individual `:bell:｜[name]` channels |
| **Claim Log** | `1445215326550622419` | Admin-only log of all claim/unclaim actions |
| **Conversion Log** | `1445214943841615922` | Admin-only log of all conversions |

> **Note:** Embed commands (`/wait-list`, `/owls-vision`, `/leaderboard`, `/panel`) are admin-only and will be executed once in whichever channel the admin chooses. No pre-configured channels required.

### 4.2 Channel Permissions Summary

```
Waitlist Channel:
├── @everyone: View ✅, Send ❌
├── Approved Owl: View ❌ (hidden after approval)
└── Full Ambassador: View ❌ (hidden)

Approved Applicants Channel (1445245329778806834):
├── @everyone: View ❌
├── Approved Owl: View ✅, Send ❌
├── Staff/Admin: View ✅, Send ✅, Manage Threads ✅
└── Full Ambassador: View ❌

Ambassador Private Channels Category (1445215664515055779):
├── @everyone: View ❌
├── Staff/Admin: View ✅, Send ✅
└── Individual ambassador channels: Only owner + staff can view

Claim Log & Conversion Log:
├── @everyone: View ❌
└── Staff/Admin: View ✅
```

---

## 5. Database Schema

### 5.1 Overview

The Ambassador System requires **PostgreSQL** tables to track:
- Approved applicants (pending full ambassador status)
- Active ambassadors
- Claims (who claimed whom)
- Conversions
- Balances and payouts
- Leaderboard data (weekly + lifetime)

### 5.2 Table Definitions

#### Table: `ambassador_applicants`
Stores users who have been approved via `/onboard` but not yet accepted as full ambassadors.

```sql
CREATE TABLE ambassador_applicants (
    id              SERIAL PRIMARY KEY,
    user_id         VARCHAR(20) NOT NULL UNIQUE,        -- Discord user ID
    username        VARCHAR(100) NOT NULL,              -- Discord username
    display_name    VARCHAR(100),                       -- Nickname/display name
    approved_by     VARCHAR(20) NOT NULL,               -- Admin who ran /onboard
    approved_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    thread_id       VARCHAR(20),                        -- Interview thread ID
    status          VARCHAR(20) DEFAULT 'pending',      -- pending, accepted, rejected
    notes           TEXT,
    
    INDEX idx_applicant_user_id (user_id),
    INDEX idx_applicant_status (status)
);
```

#### Table: `ambassadors`
Stores fully accepted ambassadors with their stats and balance.

```sql
CREATE TABLE ambassadors (
    id                  SERIAL PRIMARY KEY,
    user_id             VARCHAR(20) NOT NULL UNIQUE,    -- Discord user ID
    username            VARCHAR(100) NOT NULL,
    display_name        VARCHAR(100),
    private_channel_id  VARCHAR(20),                    -- Their :bell:｜[name] channel
    accepted_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_by         VARCHAR(20) NOT NULL,           -- Admin who accepted
    
    -- Stats
    total_claims        INT DEFAULT 0,
    total_unclaims      INT DEFAULT 0,
    total_conversions   INT DEFAULT 0,
    lifetime_conversions INT DEFAULT 0,                 -- Never resets
    current_claims      INT DEFAULT 0,                  -- Active claims (max 3)
    
    -- Balance
    balance             DECIMAL(10, 2) DEFAULT 0.00,
    total_earned        DECIMAL(10, 2) DEFAULT 0.00,
    total_paid_out      DECIMAL(10, 2) DEFAULT 0.00,
    
    -- Tier system (for future incentives)
    tier                VARCHAR(20) DEFAULT 'bronze',   -- bronze, silver, gold, platinum
    
    -- Cooldowns
    last_unclaim_at     TIMESTAMP,                      -- For 1-hour cooldown
    
    -- Status
    is_active           BOOLEAN DEFAULT TRUE,
    
    INDEX idx_ambassador_user_id (user_id),
    INDEX idx_ambassador_conversions (total_conversions DESC),
    INDEX idx_ambassador_lifetime (lifetime_conversions DESC)
);
```

#### Table: `claims`
Tracks all claim relationships between ambassadors and free users.

```sql
CREATE TABLE claims (
    id              SERIAL PRIMARY KEY,
    ambassador_id   VARCHAR(20) NOT NULL,               -- Ambassador's Discord user ID
    claimed_user_id VARCHAR(20) NOT NULL,               -- Claimed user's Discord user ID
    claimed_username VARCHAR(100),
    claimed_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status          VARCHAR(20) DEFAULT 'active',       -- active, converted, unclaimed
    converted_at    TIMESTAMP,
    unclaimed_at    TIMESTAMP,
    
    UNIQUE(ambassador_id, claimed_user_id),             -- Can't claim same person twice
    INDEX idx_claims_ambassador (ambassador_id),
    INDEX idx_claims_claimed_user (claimed_user_id),
    INDEX idx_claims_status (status)
);
```

#### Table: `conversions`
Log of all conversions for historical tracking.

```sql
CREATE TABLE conversions (
    id              SERIAL PRIMARY KEY,
    ambassador_id   VARCHAR(20) NOT NULL,
    converted_user_id VARCHAR(20) NOT NULL,
    converted_username VARCHAR(100),
    converted_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    commission      DECIMAL(10, 2) DEFAULT 0.00,        -- Amount earned
    claim_id        INT REFERENCES claims(id),
    week_number     INT,                                -- For weekly leaderboard
    year            INT,
    
    INDEX idx_conversions_ambassador (ambassador_id),
    INDEX idx_conversions_week (year, week_number)
);
```

#### Table: `payout_requests`
Tracks payout requests and their status.

```sql
CREATE TABLE payout_requests (
    id              SERIAL PRIMARY KEY,
    ambassador_id   VARCHAR(20) NOT NULL,
    amount          DECIMAL(10, 2) NOT NULL,
    payment_method  VARCHAR(100),                       -- e.g., "PayPal: email@example.com"
    thread_id       VARCHAR(20),                        -- Payout thread ID
    status          VARCHAR(20) DEFAULT 'pending',      -- pending, completed, cancelled
    requested_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at    TIMESTAMP,
    completed_by    VARCHAR(20),                        -- Admin who processed
    actual_amount   DECIMAL(10, 2),                     -- What was actually paid
    notes           TEXT,
    
    INDEX idx_payouts_ambassador (ambassador_id),
    INDEX idx_payouts_status (status)
);
```

#### Table: `balance_transactions`
Audit log of all balance changes.

```sql
CREATE TABLE balance_transactions (
    id              SERIAL PRIMARY KEY,
    ambassador_id   VARCHAR(20) NOT NULL,
    type            VARCHAR(20) NOT NULL,               -- conversion, payout, bonus, adjustment
    amount          DECIMAL(10, 2) NOT NULL,            -- Positive = credit, Negative = debit
    balance_after   DECIMAL(10, 2) NOT NULL,
    description     TEXT,
    performed_by    VARCHAR(20),                        -- Admin for manual adjustments
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_transactions_ambassador (ambassador_id),
    INDEX idx_transactions_type (type)
);
```

#### Table: `weekly_leaderboard`
Snapshot of weekly stats (resets every week).

```sql
CREATE TABLE weekly_leaderboard (
    id              SERIAL PRIMARY KEY,
    ambassador_id   VARCHAR(20) NOT NULL,
    week_number     INT NOT NULL,
    year            INT NOT NULL,
    conversions     INT DEFAULT 0,
    
    UNIQUE(ambassador_id, week_number, year),
    INDEX idx_weekly_leaderboard (year, week_number, conversions DESC)
);
```

---

## 6. Asset Requirements

### 6.1 Required Images

> ⚠️ **See top of document for complete image checklist with filenames.**

All images should be placed in `services/discordbot/assets/` and embedded via Go.

| Filename | Variable Name | Purpose |
|----------|---------------|---------|
| `waitlist.png` | `WaitlistBanner` | Header banner for `/wait-list` embed |
| `ourvision.png` | `OurVisionBanner` | Header banner for `/owls-vision` embed |
| `leaderboard.png` | `LeaderboardBanner` | Header banner for `/leaderboard` embed |
| `panel.png` | `PanelBanner` | Header banner for `/panel` embed |
| `ambassadorwelcome.png` | `AmbassadorWelcomeBanner` | Header banner for ambassador private channel welcome |
| `approved.png` | `ApprovedBanner` | Header banner shown when applicant is accepted as ambassador |
| `conversion.png` | `ConversionBanner` | Header banner for conversion success notification |
| `footer.png` | `OwlsFooter` | Footer image for ALL embeds (can reuse existing `owlsfooter.png`) |

### 6.2 Asset Embed File

Update `services/discordbot/assets/embed.go`:

```go
package assets

import _ "embed"

// Existing assets
//go:embed supportBanner.png
var SupportBanner []byte

//go:embed owlsfooter.png
var OwlsFooter []byte

// Ambassador System assets
//go:embed waitlist.png
var WaitlistBanner []byte

//go:embed ourvision.png
var OurVisionBanner []byte

//go:embed leaderboard.png
var LeaderboardBanner []byte

//go:embed panel.png
var PanelBanner []byte

//go:embed ambassadorwelcome.png
var AmbassadorWelcomeBanner []byte

//go:embed approved.png
var ApprovedBanner []byte

//go:embed conversion.png
var ConversionBanner []byte

const (
    SupportBannerFilename        = "supportBanner.png"
    OwlsFooterFilename           = "owlsfooter.png"
    WaitlistBannerFilename       = "waitlist.png"
    OurVisionBannerFilename      = "ourvision.png"
    LeaderboardBannerFilename    = "leaderboard.png"
    PanelBannerFilename          = "panel.png"
    AmbassadorWelcomeFilename    = "ambassadorwelcome.png"
    ApprovedBannerFilename       = "approved.png"
    ConversionBannerFilename     = "conversion.png"
)
```

---

## 7. Commands Reference

### 7.1 Admin Commands (Staff Only)

| Command | Server | Description |
|---------|--------|-------------|
| `/wait-list` | Ambassador | Posts the waitlist information embed with Tally button |
| `/owls-vision` | Ambassador | Posts the ambassador vision/purpose embed |
| `/onboard <user>` | Ambassador | Starts the onboarding process for an applicant |
| `/leaderboard` | Ambassador | Posts the weekly leaderboard embed |
| `/panel` | Ambassador | Posts the ambassador dashboard panel |
| `/balance-add <user> <amount>` | Ambassador | Manually add balance to an ambassador |
| `/balance-remove <user> <amount>` | Ambassador | Manually remove balance from an ambassador |
| `/ambassador-stats <user>` | Ambassador | View detailed stats for a specific ambassador |

### 7.2 Ambassador Commands

| Command | Server | Description |
|---------|--------|-------------|
| `/claim <user>` | **Main Server ONLY** | Claim a free user (max 3 active) |
| `/unclaim <user>` | **Main Server ONLY** | Release a claimed user (1hr cooldown) |
| `/my-stats` | Ambassador | View your personal stats |
| `/my-claims` | Ambassador | View your current claimed users |
| `/request-payout` | Ambassador | Request a payout (min 2 conversions, min $10) |

### 7.3 Command Server Restrictions

```
┌─────────────────────────────────────────────────────────────────┐
│                    COMMAND EXECUTION MATRIX                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  AMBASSADOR SERVER (1441909666308423874)                        │
│  ├── Admin Commands: ✅ ALL work here                           │
│  ├── /claim: ❌ Ephemeral "Wrong server" message                │
│  └── /unclaim: ❌ Ephemeral "Wrong server" message              │
│                                                                  │
│  MAIN SERVER (1431718856619200686)                              │
│  ├── Admin Commands: ❌ Do not register/show                    │
│  ├── /claim: ✅ Works here ONLY                                 │
│  └── /unclaim: ✅ Works here ONLY                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 8. Embed Specifications

All embeds use **Components V2** format with the triple-embed structure:
1. **Header:** MediaGallery with banner image
2. **Body:** Container with text content and optional buttons
3. **Footer:** MediaGallery with footer image

### 8.1 `/wait-list` Embed

```
┌─────────────────────────────────────────────────────────────────┐
│  [WAITLIST BANNER IMAGE - waitlist.png]                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  # Join the Owls Crew                                           │
│                                                                  │
│  Are you passionate about helping others succeed in their       │
│  trading journey? The Owls Capital Ambassador Program is        │
│  looking for dedicated individuals to join our elite team.      │
│                                                                  │
│  **What You'll Do:**                                            │
│  • Guide new members through their trading education            │
│  • Earn commissions for every successful conversion             │
│  • Access exclusive ambassador resources and support            │
│  • Climb the leaderboard and unlock premium rewards             │
│                                                                  │
│  **Requirements:**                                               │
│  • Active member of the Owls Capital community                  │
│  • Strong communication skills                                   │
│  • Passion for helping others succeed                           │
│  • Commitment to our community values                           │
│                                                                  │
│  Ready to apply? Click the button to start your application.    │
│                                                                  │
│                                     ┌──────────────────────┐    │
│                                     │ 📝 Apply Now         │    │
│                                     │ (Links to Tally)     │    │
│                                     └──────────────────────┘    │
├─────────────────────────────────────────────────────────────────┤
│  [FOOTER IMAGE - footer.png]                                    │
└─────────────────────────────────────────────────────────────────┘

Button Configuration:
- Style: Link Button
- Label: "Apply Now" (with 📝 emoji)
- URL: https://tally.so (placeholder — will be updated later)
- Position: RIGHT side of second embed section (use Section with Accessory)
```

### 8.2 `/owls-vision` Embed

```
┌─────────────────────────────────────────────────────────────────┐
│  [OUR VISION BANNER IMAGE - ourvision.png]                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  # Our Vision                                                    │
│                                                                  │
│  At Owls Capital, we believe success is best achieved           │
│  together. Our ambassador program isn't just about              │
│  conversions—it's about building a community of traders         │
│  who genuinely support each other's growth.                     │
│                                                                  │
│  **The Ambassador Mission:**                                     │
│                                                                  │
│  🎯 **Empower New Traders**                                     │
│  Guide newcomers through their first steps, answer              │
│  questions, and be the mentor you wish you had.                 │
│                                                                  │
│  🤝 **Build Genuine Connections**                               │
│  Success comes from authentic relationships. Focus on           │
│  helping people, not just hitting numbers.                      │
│                                                                  │
│  📈 **Lead by Example**                                         │
│  Your success story inspires others. Share your journey,        │
│  your wins, and even your lessons learned.                      │
│                                                                  │
│  🦉 **Uphold the Owls Standard**                                │
│  Represent our community with integrity. Quality over           │
│  quantity—we want the right people, not just anyone.            │
│                                                                  │
│  Together, we don't just grow numbers.                          │
│  We grow traders. We grow success stories.                      │
│  We grow the Owls family.                                       │
│                                                                  │
├─────────────────────────────────────────────────────────────────┤
│  [FOOTER IMAGE - footer.png]                                    │
└─────────────────────────────────────────────────────────────────┘

No buttons required for this embed.
```

### 8.3 `/leaderboard` Embed

```
┌─────────────────────────────────────────────────────────────────┐
│  [LEADERBOARD BANNER IMAGE - leaderboard.png]                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  # 🏆 Weekly Leaderboard                                        │
│                                                                  │
│  Top ambassadors by conversions this week.                      │
│  Resets every Sunday at midnight UTC.                           │
│                                                                  │
│  ═══════════════════════════════════════                        │
│                                                                  │
│  # 🥇 1. @TopAmbassador — 15 conversions                        │
│                                                                  │
│  ## 🥈 2. @SecondPlace — 12 conversions                         │
│                                                                  │
│  ### 🥉 3. @ThirdPlace — 10 conversions                         │
│                                                                  │
│  4. @FourthPlace — 8 conversions                                │
│  5. @FifthPlace — 7 conversions                                 │
│  6. @SixthPlace — 6 conversions                                 │
│  7. @SeventhPlace — 5 conversions                               │
│  8. @EighthPlace — 4 conversions                                │
│  9. @NinthPlace — 3 conversions                                 │
│  10. @TenthPlace — 2 conversions                                │
│                                                                  │
│  ═══════════════════════════════════════                        │
│                                                                  │
│  ┌─────────────────────────────────────┐                        │
│  │ 🏅 Lifetime Leaderboard             │                        │
│  └─────────────────────────────────────┘                        │
├─────────────────────────────────────────────────────────────────┤
│  [FOOTER IMAGE - footer.png]                                    │
└─────────────────────────────────────────────────────────────────┘

Text Sizing Rules:
- #1: Heading 1 (largest) — use "# " prefix
- #2: Heading 2 — use "## " prefix
- #3: Heading 3 — use "### " prefix
- #4-10: Regular text (no prefix)

Tie-Breaking:
- If conversions are equal, sort alphabetically by username (A→Z)

Button Configuration:
- Custom ID: "Leaderboard:ViewLifetime"
- Style: Secondary
- Label: "🏅 Lifetime Leaderboard"
- Response: Ephemeral with all-time top 10
```

### 8.4 `/panel` Embed

```
┌─────────────────────────────────────────────────────────────────┐
│  [PANEL BANNER IMAGE - panel.png]                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  # Ambassador Dashboard                                          │
│                                                                  │
│  Welcome to your command center. Use the buttons below          │
│  to manage your ambassador activities.                          │
│                                                                  │
│  ───────────────────────────────────────                        │
│                                                                  │
│  **📊 My Stats**                                                │
│  View your conversion rate, total claims, and performance       │
│  metrics.                                                       │
│                                                                  │
│  **👥 My Claims**                                               │
│  See who you're currently guiding. You can also unclaim         │
│  users here (1-hour cooldown applies).                          │
│                                                                  │
│  **💰 Request Payout**                                          │
│  Cash out your commissions. Requirements:                       │
│  • Minimum 2 conversions to unlock payouts                      │
│  • Minimum $10 balance per payout request                       │
│                                                                  │
│  ───────────────────────────────────────                        │
│                                                                  │
│  **Tier System:**                                               │
│  🥉 Bronze — Starting tier                                      │
│  🥈 Silver — 10+ lifetime conversions                           │
│  🥇 Gold — 25+ lifetime conversions                             │
│  💎 Platinum — 50+ lifetime conversions                         │
│                                                                  │
│  Higher tiers = Higher commission rates!                        │
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐          │
│  │ 📊 My Stats │  │ 👥 My Claims│  │ 💰 Request Payout│         │
│  └─────────────┘  └─────────────┘  └─────────────────┘          │
├─────────────────────────────────────────────────────────────────┤
│  [FOOTER IMAGE - footer.png]                                    │
└─────────────────────────────────────────────────────────────────┘

Button Configurations:
1. My Stats
   - Custom ID: "Panel:ViewStats"
   - Style: Primary
   - Response: Ephemeral with ambassador stats

2. My Claims
   - Custom ID: "Panel:ViewClaims"
   - Style: Primary
   - Response: Ephemeral listing current claims with unclaim buttons

3. Request Payout
   - Custom ID: "Panel:RequestPayout"
   - Style: Success (Green)
   - Response: Modal if eligible, ephemeral error if not
```

### 8.5 Ambassador Welcome Embed (Auto-generated)

Sent automatically in the ambassador's private channel after acceptance:

```
┌─────────────────────────────────────────────────────────────────┐
│  [AMBASSADOR WELCOME BANNER - ambassadorwelcome.png]            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  # 🎉 Welcome to the Team, @Username!                           │
│                                                                  │
│  Congratulations! You've officially joined the Owls Capital     │
│  Ambassador Program. This is your private channel where         │
│  you'll receive important updates and notifications.            │
│                                                                  │
│  ───────────────────────────────────────                        │
│                                                                  │
│  **Getting Started:**                                           │
│                                                                  │
│  1️⃣ Head to the Main Server to start claiming members          │
│     Use `/claim @user` on free members you want to help         │
│                                                                  │
│  2️⃣ Guide your claimed members to success                       │
│     When they upgrade to paid, you earn commission!             │
│                                                                  │
│  3️⃣ Track your progress                                        │
│     Use the Ambassador Panel to view stats and payouts          │
│                                                                  │
│  ───────────────────────────────────────                        │
│                                                                  │
│  **Your Limits:**                                               │
│  • 3 active claims at a time                                    │
│  • 1 unclaim per hour (choose wisely!)                          │
│  • Payouts unlock after 2 conversions                           │
│                                                                  │
│  **Quick Links:**                                               │
│  • 📋 Ambassador Panel — View stats & request payouts           │
│  • 🏆 Leaderboard — See where you rank                          │
│  • 📖 Our Vision — Remember why we do this                      │
│                                                                  │
│  Questions? Reach out to staff anytime. Good luck!              │
│                                                                  │
├─────────────────────────────────────────────────────────────────┤
│  [FOOTER IMAGE - footer.png]                                    │
└─────────────────────────────────────────────────────────────────┘
```

### 8.6 Interview Thread Embed (Auto-generated)

Sent in the interview thread after `/onboard`:

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│  # 📋 Ambassador Application — @Username                        │
│                                                                  │
│  Welcome to your interview! A staff member will be with         │
│  you shortly to discuss your application.                       │
│                                                                  │
│  **Application Status:** 🟡 Pending Review                      │
│                                                                  │
│  Please be patient and answer any questions from our            │
│  team honestly and thoroughly.                                  │
│                                                                  │
│  ───────────────────────────────────────                        │
│                                                                  │
│  ⚠️ **For Staff Only:**                                         │
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐                               │
│  │ ✅ Accept   │  │ ❌ Reject   │                               │
│  └─────────────┘  └─────────────┘                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

Button Configurations:
1. Accept
   - Custom ID: "Interview:Accept"
   - Style: Success (Green)
   - Response: "Are you sure?" confirmation ephemeral

2. Reject
   - Custom ID: "Interview:Reject"
   - Style: Danger (Red)
   - Response: "Are you sure?" confirmation ephemeral
```

---

## 9. Workflows

### 9.1 Waitlist → Onboarding Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          ONBOARDING WORKFLOW                              │
└──────────────────────────────────────────────────────────────────────────┘

Step 1: User Discovers Waitlist
─────────────────────────────────
• User sees /wait-list embed in public channel
• Clicks "Apply Now" button → Opens Tally form
• Completes application externally

Step 2: Admin Reviews Application (Manual)
─────────────────────────────────────────────
• Admin reviews Tally submissions manually
• If approved, admin runs: /onboard @User

Step 3: /onboard Command Execution
───────────────────────────────────
Bot performs these actions in sequence:

1. Validate user exists in Ambassador Server
2. INSERT into `ambassador_applicants` table (status: pending)
3. Grant "Approved Owl" role (1445249241743364236)
4. Update channel permissions to hide Waitlist channel from user
5. Create private thread: "Interview | [Username]"
   - Parent channel: 1445245329778806834
   - Invitable: false
   - Auto-archive: 10080 minutes (7 days)
6. Add user and staff to thread
7. Post Interview Thread Embed with Accept/Reject buttons
8. Respond to admin: "Onboarding initiated for @User"

Step 4: Interview Process
─────────────────────────
• Staff conducts interview in private thread
• Staff clicks Accept or Reject button

Step 5a: If REJECTED
────────────────────
1. Show confirmation: "Are you sure you want to reject @User?"
2. If confirmed:
   - UPDATE `ambassador_applicants` SET status = 'rejected'
   - Remove "Approved Owl" role
   - DELETE thread (not archive)
   - Optionally: DM user about rejection

Step 5b: If ACCEPTED
────────────────────
1. Show confirmation: "Are you sure you want to accept @User as an ambassador?"
2. If confirmed:
   - UPDATE `ambassador_applicants` SET status = 'accepted'
   - DELETE thread
   - Remove "Approved Owl" role
   - Grant "Full Ambassador" role (1445251690365456466)
   - Grant "Ambassador" role in Main Server (cross-server)
   - Wait 5 seconds
   - Create private channel: ":bell:｜[username]"
     - Under category: 1445215664515055779
     - Permissions: User + Staff only
   - INSERT into `ambassadors` table
   - Post Ambassador Welcome Embed in new channel
   - Log to Conversion Log: "✅ New Ambassador: @User accepted by @Admin"
```

### 9.2 Claim → Conversion Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│                           CLAIM WORKFLOW                                  │
└──────────────────────────────────────────────────────────────────────────┘

Step 1: Ambassador Executes /claim @user
─────────────────────────────────────────
Server: Main Server (1431718856619200686) ONLY

Validation Checks (in order):
┌─────────────────────────────────────────────────────────────────────────┐
│ CHECK 1: Is executor an ambassador?                                      │
│          - Must have Ambassador role in Main Server                      │
│          - Must exist in `ambassadors` table with is_active = true       │
│                                                                          │
│ CHECK 2: Is this the correct server?                                     │
│          - Guild ID must be 1431718856619200686                          │
│          - If not: "This command can only be used in the Main Server"    │
│                                                                          │
│ CHECK 3: Does ambassador have < 3 active claims?                         │
│          - Query: SELECT COUNT(*) FROM claims WHERE ambassador_id = ?    │
│            AND status = 'active'                                         │
│          - If >= 3: "You've reached your maximum of 3 active claims"     │
│                                                                          │
│ CHECK 4: Is target user claimable?                                       │
│          - NOT have Paying Member role (885910828086362132)              │
│          - NOT have Lifetime Member role (718643316786462772)            │
│          - NOT be an ambassador themselves                               │
│          - If not claimable: "This user cannot be claimed (already a     │
│            paying member or ambassador)"                                 │
│                                                                          │
│ CHECK 5: Has this user been claimed by this ambassador before?           │
│          - Query: SELECT * FROM claims WHERE ambassador_id = ? AND       │
│            claimed_user_id = ?                                           │
│          - If exists: "You have already claimed this user before"        │
│                                                                          │
│ CHECK 6: Is this user currently claimed by someone else?                 │
│          - Query: SELECT * FROM claims WHERE claimed_user_id = ? AND     │
│            status = 'active'                                             │
│          - If yes: "This user is already claimed by another ambassador"  │
└─────────────────────────────────────────────────────────────────────────┘

If all checks pass:
1. INSERT INTO claims (ambassador_id, claimed_user_id, claimed_username)
2. UPDATE ambassadors SET current_claims = current_claims + 1,
                         total_claims = total_claims + 1
3. Log to Claim Log channel (1445215326550622419):
   "📥 **CLAIM** | @Ambassador claimed @User | Total claims: X/3"
4. Respond (ephemeral): "Successfully claimed @User! Guide them well."

Step 2: Conversion Detection
────────────────────────────
Trigger: Guild Member Update event in Main Server

When a member gains the Paying Member role (885910828086362132):
1. Check if user exists in `claims` table with status = 'active'
2. If yes:
   - UPDATE claims SET status = 'converted', converted_at = NOW()
   - INSERT INTO conversions (ambassador_id, converted_user_id, ...)
   - UPDATE ambassadors SET 
       total_conversions = total_conversions + 1,
       lifetime_conversions = lifetime_conversions + 1,
       current_claims = current_claims - 1,
       balance = balance + [commission_amount]
   - INSERT INTO balance_transactions (type = 'conversion', ...)
   - INSERT/UPDATE weekly_leaderboard
   - Log to Conversion Log (1445214943841615922):
     "💰 **CONVERSION** | @User converted! Ambassador: @Ambassador | 
      Commission: $XX | Total conversions: Y"
   - Notify ambassador in private channel:
     "🎉 Great news! @User just converted! You earned $XX commission."
   - Update leaderboard embed (if already posted)
```

### 9.3 Unclaim Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          UNCLAIM WORKFLOW                                 │
└──────────────────────────────────────────────────────────────────────────┘

Command: /unclaim @user
Server: Main Server (1431718856619200686) ONLY

Validation Checks:
┌─────────────────────────────────────────────────────────────────────────┐
│ CHECK 1: Is executor an ambassador?                                      │
│                                                                          │
│ CHECK 2: Is this the correct server?                                     │
│                                                                          │
│ CHECK 3: Does ambassador have this user claimed?                         │
│          - Query: SELECT * FROM claims WHERE ambassador_id = ? AND       │
│            claimed_user_id = ? AND status = 'active'                     │
│          - If not: "You don't have this user claimed"                    │
│                                                                          │
│ CHECK 4: Is ambassador on unclaim cooldown?                              │
│          - Query: SELECT last_unclaim_at FROM ambassadors WHERE ...      │
│          - If NOW() - last_unclaim_at < 1 HOUR:                          │
│            "You can unclaim again in X minutes"                          │
└─────────────────────────────────────────────────────────────────────────┘

If all checks pass:
1. UPDATE claims SET status = 'unclaimed', unclaimed_at = NOW()
2. UPDATE ambassadors SET 
     current_claims = current_claims - 1,
     total_unclaims = total_unclaims + 1,
     last_unclaim_at = NOW()
3. Log to Claim Log (1445215326550622419):
   "📤 **UNCLAIM** | @Ambassador unclaimed @User | Claims remaining: X/3"
4. Respond (ephemeral): "You've unclaimed @User. Cooldown: 1 hour."
```

### 9.4 Payout Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│                           PAYOUT WORKFLOW                                 │
└──────────────────────────────────────────────────────────────────────────┘

Trigger: Ambassador clicks "Request Payout" button or runs /request-payout
Server: Ambassador Server ONLY

Step 1: Eligibility Check
─────────────────────────
┌─────────────────────────────────────────────────────────────────────────┐
│ CHECK 1: Does ambassador have >= 2 total conversions?                    │
│          - If not: "You need at least 2 conversions to unlock payouts"   │
│                                                                          │
│ CHECK 2: Does ambassador have >= $10 balance?                            │
│          - If not: "Minimum payout amount is $10. Your balance: $X.XX"   │
│                                                                          │
│ CHECK 3: Does ambassador have a pending payout request?                  │
│          - If yes: "You already have a pending payout request"           │
└─────────────────────────────────────────────────────────────────────────┘

Step 2: Show Modal
──────────────────
Modal Title: "Request Payout"

Fields:
┌─────────────────────────────────────────────────────────────────────────┐
│ Field 1: Payout Amount                                                   │
│ - Label: "How much would you like to withdraw? (Your balance: $XX.XX)"   │
│ - Placeholder: "Enter amount (e.g., 25.00)"                              │
│ - Style: Short                                                           │
│ - Required: true                                                         │
│ - Validation: Must be number, <= balance, >= 10                          │
│                                                                          │
│ Field 2: Payment Method                                                  │
│ - Label: "Preferred Payment Method"                                      │
│ - Placeholder: "PayPal email address or other payment details"           │
│ - Style: Paragraph                                                       │
│ - Required: true                                                         │
└─────────────────────────────────────────────────────────────────────────┘

Step 3: Process Request
───────────────────────
1. Validate modal input
2. INSERT INTO payout_requests (ambassador_id, amount, payment_method, ...)
3. Create private thread in Panel channel:
   - Name: "[PAYOUT] @Username — $XX.XX"
   - Add ambassador and staff
4. Post payout request embed in thread:

   ┌───────────────────────────────────────────────────────────────────┐
   │ # 💰 Payout Request                                               │
   │                                                                   │
   │ **Ambassador:** @Username                                         │
   │ **Amount Requested:** $XX.XX                                      │
   │ **Current Balance:** $XX.XX                                       │
   │ **Payment Method:** PayPal: email@example.com                     │
   │                                                                   │
   │ **Ambassador Stats:**                                             │
   │ • Total Conversions: XX                                           │
   │ • Lifetime Earnings: $XX.XX                                       │
   │ • Total Paid Out: $XX.XX                                          │
   │                                                                   │
   │ ───────────────────────────────────────                          │
   │                                                                   │
   │ ⚠️ Staff: Send payment, then close this ticket.                   │
   │                                                                   │
   │ ┌───────────────────┐                                            │
   │ │ ✅ Complete Payout │                                            │
   │ └───────────────────┘                                            │
   └───────────────────────────────────────────────────────────────────┘

5. Respond (ephemeral): "Payout request submitted! Check your thread."

Step 4: Admin Completes Payout
──────────────────────────────
Admin clicks "Complete Payout" button:

1. Show modal:
   - Title: "Complete Payout"
   - Field: "Actual amount paid to ambassador"
   - Placeholder: "Enter the exact amount you sent"

2. On submit:
   - UPDATE payout_requests SET status = 'completed', 
       completed_at = NOW(), actual_amount = [input], completed_by = [admin]
   - UPDATE ambassadors SET 
       balance = balance - [actual_amount],
       total_paid_out = total_paid_out + [actual_amount]
   - INSERT INTO balance_transactions (type = 'payout', amount = -[actual], ...)
   - Delete thread
   - Notify ambassador in private channel:
     "💸 Payout complete! $XX.XX has been sent to your PayPal."
   - Log to ambassador's private channel:
     "📝 Transaction: -$XX.XX (Payout) | New balance: $XX.XX"
```

---

## 10. Business Rules

### 10.1 Claiming Rules

| Rule | Description |
|------|-------------|
| **Max Active Claims** | 3 users at a time |
| **Duplicate Claims** | Cannot claim the same user twice ever |
| **Claimable Users** | Free members only (no paying, lifetime, or ambassador roles) |
| **Claim Location** | Main Server (1431718856619200686) ONLY |
| **Claim Persistence** | Claims remain until conversion, unclaim, or target becomes unclaimable |

### 10.2 Unclaim Rules

| Rule | Description |
|------|-------------|
| **Cooldown** | 1 hour between unclaims |
| **Unclaim Location** | Main Server (1431718856619200686) ONLY |
| **Effect** | Frees up claim slot, user becomes claimable by others |
| **Tracking** | Unclaims counted in ambassador stats |

### 10.3 Conversion Rules

| Rule | Description |
|------|-------------|
| **Trigger** | Claimed user gains Paying Member role (885910828086362132) |
| **Attribution** | Full credit to claiming ambassador |
| **Commission** | Amount determined by tier (define in config) |
| **Weekly Reset** | Weekly conversion count resets Sunday midnight UTC |
| **Lifetime** | Lifetime conversions never reset |

### 10.4 Payout Rules

| Rule | Description |
|------|-------------|
| **Unlock Requirement** | Minimum 2 conversions (ever) |
| **Minimum Payout** | $10.00 |
| **Payment Methods** | PayPal (primary), others as agreed |
| **Balance Deduction** | Upon admin confirmation, not request |
| **Pending Limit** | 1 pending request at a time |

### 10.5 Tier System (Suggested)

| Tier | Requirement | Commission Rate |
|------|-------------|-----------------|
| 🥉 Bronze | Default | Base rate |
| 🥈 Silver | 10+ lifetime conversions | +10% |
| 🥇 Gold | 25+ lifetime conversions | +20% |
| 💎 Platinum | 50+ lifetime conversions | +30% |

---

## 11. Logging System

### 11.1 Claim Log (1445215326550622419)

Events logged:
- Claims: `📥 **CLAIM** | @Ambassador claimed @User | Claims: X/3`
- Unclaims: `📤 **UNCLAIM** | @Ambassador unclaimed @User | Claims: X/3`

Format: Embed with fields
```
Title: 📥 Claim Logged
Color: Green (#00FF00)
Fields:
- Ambassador: @User (ID: 123456)
- Claimed User: @User (ID: 789012)  
- Claims Active: 2/3
- Timestamp: <t:UNIX:F>
```

### 11.2 Conversion Log (1445214943841615922)

Events logged:
- Conversions: `💰 **CONVERSION** | @User converted by @Ambassador | +$XX`
- New Ambassadors: `🎉 **NEW AMBASSADOR** | @User accepted by @Admin`

Format: Embed with fields
```
Title: 💰 Conversion Recorded
Color: Gold (#FFD700)
Fields:
- Ambassador: @User (ID: 123456)
- Converted User: @User (ID: 789012)
- Commission Earned: $XX.XX
- Total Conversions: XX
- Timestamp: <t:UNIX:F>
```

### 11.3 Ambassador Private Channel Logs

Events logged:
- Balance changes (conversions, payouts, bonuses, adjustments)
- Payout completions
- Tier upgrades

Format: Simple text messages with emoji indicators

---

## 12. Permission Matrix

### 12.1 Command Permissions

| Command | Admin | Ambassador | Server |
|---------|-------|------------|--------|
| `/wait-list` | ✅ | ❌ | Ambassador |
| `/owls-vision` | ✅ | ❌ | Ambassador |
| `/onboard` | ✅ | ❌ | Ambassador |
| `/leaderboard` | ✅ | ❌ | Ambassador |
| `/panel` | ✅ | ❌ | Ambassador |
| `/balance-add` | ✅ | ❌ | Ambassador |
| `/balance-remove` | ✅ | ❌ | Ambassador |
| `/ambassador-stats` | ✅ | ❌ | Ambassador |
| `/claim` | ❌ | ✅ | **Main** |
| `/unclaim` | ❌ | ✅ | **Main** |
| `/my-stats` | ❌ | ✅ | Ambassador |
| `/my-claims` | ❌ | ✅ | Ambassador |
| `/request-payout` | ❌ | ✅ | Ambassador |

### 12.2 Button Permissions

| Button | Who Can Click |
|--------|---------------|
| Apply Now (Waitlist) | Anyone |
| Lifetime Leaderboard | Anyone |
| My Stats | Ambassadors only |
| My Claims | Ambassadors only |
| Request Payout | Ambassadors only |
| Accept/Reject (Interview) | Staff only |
| Complete Payout | Staff only |

### 12.3 Channel Access

| Channel | Everyone | Approved Owl | Ambassador | Staff |
|---------|----------|--------------|------------|-------|
| Waitlist | ✅ View | ❌ Hidden | ❌ Hidden | ✅ Full |
| Approved Applicants | ❌ | ✅ View | ❌ | ✅ Full |
| General Amb Channels | ❌ | ❌ | ✅ Full | ✅ Full |
| Private Amb Channels | ❌ | ❌ | Owner Only | ✅ Full |
| Claim Log | ❌ | ❌ | ❌ | ✅ View |
| Conversion Log | ❌ | ❌ | ❌ | ✅ View |

---

## 13. Pre-Implementation Checklist

### 13.1 Discord Server Setup

#### Ambassador Server (1441909666308423874)

- [ ] Verify Staff/Admin role exists: `1441909798852366499`
- [ ] Verify Approved Owl role exists: `1445249241743364236`
- [ ] Verify Full Ambassador role exists: `1445251690365456466`
- [ ] Verify Approved Applicants channel exists: `1445245329778806834`
- [ ] Verify Private Channels category exists: `1445215664515055779`
- [ ] Verify Claim Log channel exists: `1445215326550622419`
- [ ] Verify Conversion Log channel exists: `1445214943841615922`
- [ ] Ensure bot role is **highest** in hierarchy
- [ ] Bot has Administrator permission (or: Manage Channels, Manage Roles, Send Messages, Create Threads, Manage Threads, Embed Links, Attach Files)

#### Main Server (1431718856619200686)

- [ ] Verify Default Member role: `718643370301325404`
- [ ] Verify Paying Member role: `885910828086362132`
- [ ] Verify Lifetime Member role: `718643316786462772`
- [ ] Create Ambassador role (for cross-server identification)
- [ ] Record new Ambassador role ID: `________________`
- [ ] Ensure bot role is **highest** in hierarchy
- [ ] Bot has permissions: Manage Roles, View Channels, Read Message History

### 13.2 Assets

- [ ] `waitlist.png` — Waitlist header banner (1200×400 recommended)
- [ ] `ourvision.png` — Vision header banner
- [ ] `leaderboard.png` — Leaderboard header banner
- [ ] `panel.png` — Panel header banner
- [ ] `ambassadorwelcome.png` — Welcome header banner
- [ ] `approved.png` — Approved/accepted applicant banner
- [ ] `conversion.png` — Conversion success notification banner
- [ ] `footer.png` — Footer image (can reuse existing `owlsfooter.png`)
- [ ] Update `assets/embed.go` with new embeds

### 13.3 Database

- [ ] PostgreSQL server running and accessible
- [ ] Create `ambassador_applicants` table
- [ ] Create `ambassadors` table
- [ ] Create `claims` table
- [ ] Create `conversions` table
- [ ] Create `payout_requests` table
- [ ] Create `balance_transactions` table
- [ ] Create `weekly_leaderboard` table
- [ ] Create indexes as specified
- [ ] Test database connectivity from bot

### 13.4 External Services

- [ ] Tally form created and URL updated (placeholder: `https://tally.so`)
- [ ] PayPal or payment method configured for payouts

### 13.5 Configuration Values

Add to environment/config:

```env
# Ambassador System Configuration
AMBASSADOR_SERVER_ID=1441909666308423874
MAIN_SERVER_ID=1431718856619200686

# Ambassador Server Roles
STAFF_ROLE_ID=1441909798852366499
APPROVED_OWL_ROLE_ID=1445249241743364236
FULL_AMBASSADOR_ROLE_ID=1445251690365456466

# Main Server Roles  
MAIN_AMBASSADOR_ROLE_ID=________________  # Fill after creating
PAYING_MEMBER_ROLE_ID=885910828086362132
LIFETIME_MEMBER_ROLE_ID=718643316786462772
DEFAULT_MEMBER_ROLE_ID=718643370301325404

# Ambassador Server Channels
APPROVED_APPLICANTS_CHANNEL_ID=1445245329778806834
AMBASSADOR_CATEGORY_ID=1445215664515055779
CLAIM_LOG_CHANNEL_ID=1445215326550622419
CONVERSION_LOG_CHANNEL_ID=1445214943841615922

# Business Rules
MAX_ACTIVE_CLAIMS=3
UNCLAIM_COOLDOWN_MINUTES=60
MIN_CONVERSIONS_FOR_PAYOUT=2
MIN_PAYOUT_AMOUNT=10.00

# Commission Rates (per tier)
BRONZE_COMMISSION=10.00
SILVER_COMMISSION=11.00
GOLD_COMMISSION=12.00
PLATINUM_COMMISSION=13.00

# Tier Thresholds
SILVER_THRESHOLD=10
GOLD_THRESHOLD=25
PLATINUM_THRESHOLD=50

# External URLs
TALLY_FORM_URL=https://tally.so  # Placeholder — update with actual form URL later
```

---

## 14. Technical Implementation Notes

### 14.1 File Structure

```
services/discordbot/internal/features/ambassador/
├── ambassador.go           # Main feature registration
├── config.go               # Configuration constants and structs
├── database.go             # Database operations
├── embeds.go               # All embed builders
├── handlers_admin.go       # Admin command handlers
├── handlers_ambassador.go  # Ambassador command handlers
├── handlers_buttons.go     # Button interaction handlers
├── handlers_modals.go      # Modal submission handlers
├── workflows.go            # Complex multi-step workflows
├── leaderboard.go          # Leaderboard logic and updates
├── conversion_tracker.go   # Role change event listener
└── utils.go                # Helper functions
```

### 14.2 Intents Required

```go
func (f *AmbassadorFeature) Intents() discordgo.Intent {
    return discordgo.IntentGuilds |
           discordgo.IntentGuildMembers |    // For role change detection
           discordgo.IntentGuildMessages
}
```

### 14.3 Event Handlers Required

```go
// Guild Member Update - for conversion detection
func (f *AmbassadorFeature) OnGuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
    // Check if user gained paying role
    // If claimed, trigger conversion workflow
}
```

### 14.4 Rate Limit Considerations

- Use bulk operations where possible
- Cache ambassador data in memory with TTL
- Implement exponential backoff for API calls
- Queue non-critical log messages
- Use database transactions for multi-step operations

### 14.5 Error Handling

- All database operations should use transactions
- Rollback on any failure in multi-step workflows
- Log all errors with context
- Provide user-friendly error messages
- Never expose internal errors to users

### 14.6 Weekly Reset Mechanism

Implement a scheduled task (cron or ticker) that runs every Sunday at midnight UTC:

```go
func (f *AmbassadorFeature) WeeklyReset() {
    // 1. Archive current week's leaderboard
    // 2. Reset weekly conversion counts in ambassadors table
    // 3. Update leaderboard embed message
    // 4. Optionally: Announce weekly winners
}
```

---

## Appendix A: Quick Reference — All IDs

### Servers
| Name | ID |
|------|-----|
| Main Server | `1431718856619200686` |
| Ambassador Server | `1441909666308423874` |

### Ambassador Server Roles
| Name | ID |
|------|-----|
| Staff/Admin | `1441909798852366499` |
| Approved Owl | `1445249241743364236` |
| Full Ambassador | `1445251690365456466` |

### Main Server Roles
| Name | ID |
|------|-----|
| Default Member | `718643370301325404` |
| Paying Member | `885910828086362132` |
| Lifetime Member | `718643316786462772` |
| Ambassador | _TO BE CREATED_ |

### Ambassador Server Channels
| Name | ID |
|------|-----|
| Approved Applicants | `1445245329778806834` |
| Private Channels Category | `1445215664515055779` |
| Claim Log | `1445215326550622419` |
| Conversion Log | `1445214943841615922` |

---

## Appendix B: Sample Database Queries

### Get Ambassador Stats
```sql
SELECT 
    a.user_id,
    a.username,
    a.total_claims,
    a.total_conversions,
    a.lifetime_conversions,
    a.current_claims,
    a.balance,
    a.tier,
    CASE WHEN a.total_claims > 0 
         THEN ROUND(a.total_conversions::decimal / a.total_claims * 100, 2)
         ELSE 0 
    END as conversion_rate
FROM ambassadors a
WHERE a.user_id = $1 AND a.is_active = true;
```

### Get Weekly Leaderboard
```sql
SELECT 
    a.user_id,
    a.username,
    COALESCE(w.conversions, 0) as weekly_conversions
FROM ambassadors a
LEFT JOIN weekly_leaderboard w ON a.user_id = w.ambassador_id 
    AND w.week_number = EXTRACT(WEEK FROM CURRENT_DATE)
    AND w.year = EXTRACT(YEAR FROM CURRENT_DATE)
WHERE a.is_active = true
ORDER BY weekly_conversions DESC, a.username ASC
LIMIT 10;
```

### Get Lifetime Leaderboard
```sql
SELECT 
    user_id,
    username,
    lifetime_conversions
FROM ambassadors
WHERE is_active = true
ORDER BY lifetime_conversions DESC, username ASC
LIMIT 10;
```

### Check Claim Eligibility
```sql
-- Check if ambassador can claim more users
SELECT current_claims < 3 as can_claim
FROM ambassadors
WHERE user_id = $1 AND is_active = true;

-- Check if user already claimed by this ambassador
SELECT EXISTS(
    SELECT 1 FROM claims 
    WHERE ambassador_id = $1 AND claimed_user_id = $2
) as already_claimed;

-- Check if user currently claimed by anyone
SELECT EXISTS(
    SELECT 1 FROM claims 
    WHERE claimed_user_id = $1 AND status = 'active'
) as currently_claimed;
```

### Check Unclaim Cooldown
```sql
SELECT 
    CASE WHEN last_unclaim_at IS NULL THEN true
         WHEN NOW() - last_unclaim_at > INTERVAL '1 hour' THEN true
         ELSE false
    END as can_unclaim,
    CASE WHEN last_unclaim_at IS NOT NULL 
              AND NOW() - last_unclaim_at <= INTERVAL '1 hour'
         THEN EXTRACT(EPOCH FROM (last_unclaim_at + INTERVAL '1 hour' - NOW())) / 60
         ELSE 0
    END as minutes_remaining
FROM ambassadors
WHERE user_id = $1;
```

---

## Appendix C: Embed Color Palette

| Purpose | Hex Color | Usage |
|---------|-----------|-------|
| Primary/Brand | `#FFB800` | Main embeds, success states |
| Success | `#00FF00` | Claim logs, approvals |
| Warning | `#FFA500` | Pending states, cooldowns |
| Danger | `#FF0000` | Rejections, errors |
| Info | `#0099FF` | Stats, information |
| Conversion | `#FFD700` | Conversion notifications |

---

> **Document End**
> 
> This specification should provide complete guidance for implementing the Ambassador System.
> For questions or clarifications, consult with the project lead.

