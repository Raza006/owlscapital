# Freemium Module System — Complete Specification

> **Purpose:** Build and manage educational modules where free members can preview all content but only access the first lesson, while premium members unlock everything.

---

## 📋 TABLE OF CONTENTS

1. [System Overview](#-system-overview)
2. [Role Configuration](#-role-configuration)
3. [Database Schema](#-database-schema)
4. [Module Structure](#-module-structure)
5. [Commands Reference](#-commands-reference)
6. [Flowcharts](#-flowcharts)
7. [Components V2 Layouts](#-components-v2-layouts)
8. [Setup Checklist](#-setup-checklist)
9. [Assets Required](#-assets-required)

---

## 🎯 SYSTEM OVERVIEW

### The Freemium Concept

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           FREEMIUM PHILOSOPHY                                │
│                                                                              │
│   "See Everything, Access What You've Earned"                               │
│                                                                              │
│   Free Members:     Can SEE all lessons → Can ACCESS only Lesson 1          │
│   Premium Members:  Can SEE all lessons → Can ACCESS all lessons            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Architecture

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              SYSTEM ARCHITECTURE                              │
└──────────────────────────────────────────────────────────────────────────────┘

                    ┌─────────────────────────────┐
                    │      ADMIN/STAFF            │
                    │                             │
                    │  Uses /module commands      │
                    └──────────────┬──────────────┘
                                   │
                                   ▼
                    ┌─────────────────────────────┐
                    │      MODULE BUILDER         │
                    │                             │
                    │  • /module create           │
                    │  • /module edit             │
                    │  • /module lesson-add       │
                    │  • /module lesson-remove    │
                    │  • /module publish          │
                    │  • /module list             │
                    │  • /module delete           │
                    └──────────────┬──────────────┘
                                   │
                                   ▼
                    ┌─────────────────────────────┐
                    │      POSTGRESQL DATABASE    │
                    │                             │
                    │  • modules table            │
                    │  • lessons table            │
                    │  • lesson_access_log        │
                    └──────────────┬──────────────┘
                                   │
                                   ▼
                    ┌─────────────────────────────┐
                    │      FORUM CHANNEL          │
                    │   (1447780920940695655)     │
                    │                             │
                    │  Each module = 1 Forum Post │
                    │  with Components V2 embed   │
                    └──────────────┬──────────────┘
                                   │
                                   ▼
              ┌────────────────────┴────────────────────┐
              │                                         │
              ▼                                         ▼
┌─────────────────────────┐             ┌─────────────────────────┐
│     FREE MEMBER         │             │    PREMIUM MEMBER       │
│                         │             │                         │
│  Selects from dropdown  │             │  Selects from dropdown  │
│  → Only Lesson 1 works  │             │  → ALL lessons work     │
│  → Others show upgrade  │             │  → Full content access  │
│    prompt               │             │                         │
└─────────────────────────┘             └─────────────────────────┘
```

---

## 🔐 ROLE CONFIGURATION

### Required Roles

| Role | ID | Access Level | Description |
|------|-----|--------------|-------------|
| **Lifetime Premium** | `718643316786462772` | FULL | Permanent access to all lessons |
| **Premium** | `885910828086362132` | FULL | Subscription-based access to all lessons |
| **Free/Default** | `718643370301325404` | LIMITED | First lesson only per module |

### Permission Logic

```go
// Pseudo-code for access check
func canAccessLesson(member, lesson) bool {
    // Premium or Lifetime = full access
    if hasRole(member, LIFETIME_ROLE) || hasRole(member, PREMIUM_ROLE) {
        return true
    }
    
    // Free members can only access lesson with order_index = 1
    if lesson.OrderIndex == 1 {
        return true
    }
    
    // All other cases = denied
    return false
}
```

### Bot Permissions Required

```
┌─────────────────────────────────────────────────────────────────┐
│                     BOT PERMISSIONS                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  REQUIRED:                                                       │
│  ✅ Manage Channels (for forum post creation)                    │
│  ✅ Send Messages                                                │
│  ✅ Send Messages in Threads                                     │
│  ✅ Create Public Threads (forum posts)                          │
│  ✅ Embed Links                                                  │
│  ✅ Attach Files                                                 │
│  ✅ Use External Emojis                                          │
│  ✅ Read Message History                                         │
│  ✅ Add Reactions                                                │
│                                                                  │
│  RECOMMENDED:                                                    │
│  ✅ Manage Messages (to pin/edit module embeds)                  │
│  ✅ Manage Threads (to manage forum posts)                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🗄️ DATABASE SCHEMA

### Table: `freemium_modules`

```sql
CREATE TABLE freemium_modules (
    -- Identity
    id                  SERIAL PRIMARY KEY,
    module_id           VARCHAR(50) UNIQUE NOT NULL,      -- Unique slug: "trading-basics"
    
    -- Display
    title               VARCHAR(100) NOT NULL,            -- "Trading Basics 101"
    description         TEXT NOT NULL,                    -- Module description (markdown supported)
    banner_url          VARCHAR(500),                     -- Custom banner image URL (optional)
    
    -- Discord References
    forum_channel_id    VARCHAR(20) NOT NULL,             -- 1447780920940695655
    forum_post_id       VARCHAR(20),                      -- Created forum post/thread ID
    embed_message_id    VARCHAR(20),                      -- Message ID of the embed in the post
    
    -- Metadata
    created_by          VARCHAR(20) NOT NULL,             -- Admin who created it
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Status
    is_published        BOOLEAN DEFAULT FALSE,            -- Draft vs Published
    is_archived         BOOLEAN DEFAULT FALSE             -- Soft delete
);

-- Index for quick lookups
CREATE INDEX idx_modules_forum_post ON freemium_modules(forum_post_id);
CREATE INDEX idx_modules_published ON freemium_modules(is_published) WHERE is_published = TRUE;
```

### Table: `freemium_lessons`

```sql
CREATE TABLE freemium_lessons (
    -- Identity
    id                  SERIAL PRIMARY KEY,
    module_id           VARCHAR(50) REFERENCES freemium_modules(module_id) ON DELETE CASCADE,
    lesson_id           VARCHAR(50) NOT NULL,             -- Unique within module: "lesson-1"
    
    -- Display
    title               VARCHAR(100) NOT NULL,            -- "Introduction to Charts"
    description         TEXT,                             -- Brief description for dropdown
    content             TEXT NOT NULL,                    -- Full lesson content (markdown)
    content_image_url   VARCHAR(500),                     -- Optional image for lesson content
    
    -- Organization
    order_index         INT NOT NULL,                     -- 1, 2, 3... (1 = free lesson)
    
    -- Access Control
    is_free             BOOLEAN GENERATED ALWAYS AS (order_index = 1) STORED,
    
    -- Metadata
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    UNIQUE(module_id, lesson_id),
    UNIQUE(module_id, order_index)
);

-- Index for ordering
CREATE INDEX idx_lessons_order ON freemium_lessons(module_id, order_index);
```

### Table: `freemium_access_log`

```sql
CREATE TABLE freemium_access_log (
    -- Identity
    id                  SERIAL PRIMARY KEY,
    
    -- References
    user_id             VARCHAR(20) NOT NULL,
    module_id           VARCHAR(50) NOT NULL,
    lesson_id           VARCHAR(50) NOT NULL,
    
    -- Access Details
    access_granted      BOOLEAN NOT NULL,                 -- TRUE = viewed, FALSE = blocked
    user_role_type      VARCHAR(20) NOT NULL,             -- 'free', 'premium', 'lifetime'
    
    -- Timestamp
    accessed_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for analytics
CREATE INDEX idx_access_user ON freemium_access_log(user_id);
CREATE INDEX idx_access_module ON freemium_access_log(module_id);
CREATE INDEX idx_access_time ON freemium_access_log(accessed_at);
```

### Table: `freemium_settings`

```sql
CREATE TABLE freemium_settings (
    -- Identity
    id                  SERIAL PRIMARY KEY,
    setting_key         VARCHAR(50) UNIQUE NOT NULL,
    
    -- Value
    setting_value       TEXT NOT NULL,
    
    -- Metadata
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by          VARCHAR(20)
);

-- Default settings
INSERT INTO freemium_settings (setting_key, setting_value) VALUES
    ('forum_channel_id', '1447780920940695655'),
    ('default_footer_url', 'attachment://footer.png'),
    ('lifetime_role_id', '718643316786462772'),
    ('premium_role_id', '885910828086362132'),
    ('free_role_id', '718643370301325404'),
    ('upgrade_message', 'Upgrade to Premium to unlock all lessons!'),
    ('upgrade_url', 'https://your-upgrade-link.com');
```

---

## 📦 MODULE STRUCTURE

### What is a Module?

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              MODULE ANATOMY                                  │
└─────────────────────────────────────────────────────────────────────────────┘

MODULE: "Trading Basics 101"
│
├─── module_id: "trading-basics"
├─── title: "Trading Basics 101"
├─── description: "Learn the fundamentals of trading..."
├─── banner_url: "https://cdn.../trading-banner.png"  (or attachment)
│
└─── LESSONS (Ordered)
     │
     ├─── [1] "Introduction" ──────────────── ✅ FREE (Everyone)
     │    └── content: "Welcome to trading..."
     │
     ├─── [2] "Understanding Charts" ──────── 🔒 PREMIUM ONLY
     │    └── content: "Charts are the..."
     │
     ├─── [3] "Risk Management" ───────────── 🔒 PREMIUM ONLY
     │    └── content: "Never risk more..."
     │
     ├─── [4] "Entry Strategies" ──────────── 🔒 PREMIUM ONLY
     │    └── content: "When to enter..."
     │
     └─── [5] "Advanced Patterns" ─────────── 🔒 PREMIUM ONLY
          └── content: "Complex patterns..."
```

### Forum Post Layout

When a module is published, it creates a **Forum Post** in the designated forum channel:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  FORUM: Educational Modules                                                  │
│  Channel ID: 1447780920940695655                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  📌 Trading Basics 101                         [Posted by Bot]              │
│  📌 Technical Analysis Mastery                 [Posted by Bot]              │
│  📌 Risk Management Fundamentals               [Posted by Bot]              │
│  📌 Psychology of Trading                      [Posted by Bot]              │
│  📌 Options Trading 101                        [Posted by Bot]              │
│  📌 Crypto Fundamentals                        [Posted by Bot]              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 🔧 COMMANDS REFERENCE

### Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           COMMAND STRUCTURE                                  │
└─────────────────────────────────────────────────────────────────────────────┘

/module
├── create      Create a new module (draft)
├── edit        Edit an existing module
├── delete      Delete a module permanently
├── list        List all modules (drafts & published)
├── publish     Publish a draft to the forum
├── unpublish   Remove from forum (keep in DB)
├── lesson
│   ├── add       Add a lesson to a module
│   ├── edit      Edit an existing lesson
│   ├── remove    Remove a lesson from a module
│   └── reorder   Change lesson order
└── settings    Configure system settings (admin)
```

---

### `/module create`

**Purpose:** Create a new module in draft state.

**Options:**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `name` | String | ✅ | Module title (e.g., "Trading Basics 101") |
| `description` | String | ✅ | Module description (supports markdown) |
| `banner` | Attachment | ❌ | Custom banner image (optional) |

**Flow:**

```
/module create name:"Trading Basics 101" description:"Learn the fundamentals..."
      │
      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              VALIDATION                                      │
└─────────────────────────────────────────────────────────────────────────────┘
      │
      ├──── CHECK: Is user Staff/Admin?
      │     ├─── NO ──▶ ❌ "Only administrators can create modules."
      │     └─── YES ─▶ CONTINUE
      │
      ├──── CHECK: Does module with same name exist?
      │     ├─── YES ─▶ ❌ "A module with this name already exists."
      │     └─── NO ──▶ CONTINUE
      │
      └──── GENERATE: module_id from title (slug)
            │
            │  "Trading Basics 101" → "trading-basics-101"
            │
            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              EXECUTION                                       │
└─────────────────────────────────────────────────────────────────────────────┘
      │
      │  STEP 1: Upload banner (if provided) or use default
      │
      ▼
      │  STEP 2: Insert into database
      │
      │  INSERT INTO freemium_modules (
      │      module_id, title, description, banner_url,
      │      forum_channel_id, created_by, is_published
      │  ) VALUES (
      │      'trading-basics-101', 'Trading Basics 101', '...',
      │      banner_url, '1447780920940695655', admin_id, FALSE
      │  )
      │
      ▼
      │  STEP 3: Respond with success embed
      │
      │  ┌─────────────────────────────────────────────────────────────┐
      │  │ ✅ **Module Created (Draft)**                               │
      │  │                                                             │
      │  │ **Title:** Trading Basics 101                               │
      │  │ **ID:** `trading-basics-101`                                │
      │  │ **Status:** 📝 Draft                                        │
      │  │                                                             │
      │  │ **Next Steps:**                                             │
      │  │ 1. Add lessons with `/module lesson add`                    │
      │  │ 2. Preview with `/module preview`                           │
      │  │ 3. Publish with `/module publish`                           │
      │  │                                                             │
      │  └─────────────────────────────────────────────────────────────┘
      │
      └─────▶ END
```

---

### `/module lesson add`

**Purpose:** Add a lesson to an existing module.

**Options:**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `module` | String (Autocomplete) | ✅ | Select module to add lesson to |
| `title` | String | ✅ | Lesson title |
| `content` | String | ✅ | Lesson content (markdown supported) |
| `image` | Attachment | ❌ | Lesson image (optional) |
| `position` | Integer | ❌ | Position (default: append to end) |

**Flow:**

```
/module lesson add module:"trading-basics-101" title:"Introduction" content:"Welcome..."
      │
      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              VALIDATION                                      │
└─────────────────────────────────────────────────────────────────────────────┘
      │
      ├──── CHECK: Is user Staff/Admin?
      │     ├─── NO ──▶ ❌ "Only administrators can manage modules."
      │     └─── YES ─▶ CONTINUE
      │
      ├──── CHECK: Does module exist?
      │     ├─── NO ──▶ ❌ "Module not found."
      │     └─── YES ─▶ CONTINUE
      │
      └──── CHECK: Is module unpublished (editable)?
            │     (Published modules can still add lessons, but require re-publish)
            │
            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              EXECUTION                                       │
└─────────────────────────────────────────────────────────────────────────────┘
      │
      │  STEP 1: Calculate order_index
      │
      │  SELECT COALESCE(MAX(order_index), 0) + 1 AS next_index
      │  FROM freemium_lessons WHERE module_id = ?
      │
      ▼
      │  STEP 2: Generate lesson_id
      │
      │  "Introduction" → "introduction"
      │  (or "lesson-{order_index}" if collision)
      │
      ▼
      │  STEP 3: Insert lesson
      │
      │  INSERT INTO freemium_lessons (
      │      module_id, lesson_id, title, content,
      │      content_image_url, order_index
      │  ) VALUES (...)
      │
      ▼
      │  STEP 4: Respond
      │
      │  ┌─────────────────────────────────────────────────────────────┐
      │  │ ✅ **Lesson Added**                                         │
      │  │                                                             │
      │  │ **Module:** Trading Basics 101                              │
      │  │ **Lesson:** Introduction                                    │
      │  │ **Position:** #1 (FREE lesson)                              │
      │  │                                                             │
      │  │ ℹ️ Position #1 is always free for all members.              │
      │  │                                                             │
      │  │ **Current Lessons:**                                        │
      │  │ 1. ✅ Introduction (FREE)                                   │
      │  │                                                             │
      │  └─────────────────────────────────────────────────────────────┘
      │
      └─────▶ END
```

---

### `/module publish`

**Purpose:** Publish a draft module to the forum channel.

**Options:**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `module` | String (Autocomplete) | ✅ | Select module to publish |

**Flow:**

```
/module publish module:"trading-basics-101"
      │
      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              VALIDATION                                      │
└─────────────────────────────────────────────────────────────────────────────┘
      │
      ├──── CHECK: Module exists?
      │     ├─── NO ──▶ ❌ "Module not found."
      │     └─── YES ─▶ CONTINUE
      │
      ├──── CHECK: Module has at least 1 lesson?
      │     ├─── NO ──▶ ❌ "Module must have at least 1 lesson to publish."
      │     └─── YES ─▶ CONTINUE
      │
      └──── CHECK: Module already published?
            ├─── YES ─▶ UPDATE existing forum post
            └─── NO ──▶ CREATE new forum post
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         CREATE FORUM POST                                    │
└─────────────────────────────────────────────────────────────────────────────┘
      │
      │  STEP 1: Build Components V2 Embed
      │  (See Components V2 section below)
      │
      ▼
      │  STEP 2: Create Forum Post (Thread)
      │
      │  s.ForumThreadStart(forum_channel_id, &discordgo.ThreadStart{
      │      Name:                module.Title,
      │      AutoArchiveDuration: 10080,  // 7 days
      │  }, &discordgo.MessageSend{
      │      Components: componentsV2Embed,
      │      Files:      [banner, footer],
      │  })
      │
      ▼
      │  STEP 3: Store post ID in database
      │
      │  UPDATE freemium_modules SET
      │      forum_post_id = ?,
      │      embed_message_id = ?,
      │      is_published = TRUE,
      │      updated_at = NOW()
      │  WHERE module_id = ?
      │
      ▼
      │  STEP 4: Respond to admin
      │
      │  ┌─────────────────────────────────────────────────────────────┐
      │  │ ✅ **Module Published!**                                    │
      │  │                                                             │
      │  │ **Module:** Trading Basics 101                              │
      │  │ **Forum Post:** [Click to View](link)                       │
      │  │ **Lessons:** 5                                              │
      │  │                                                             │
      │  │ Members can now access this module!                         │
      │  │                                                             │
      │  └─────────────────────────────────────────────────────────────┘
      │
      └─────▶ END
```

---

### `/module edit`

**Purpose:** Edit module title, description, or banner.

**Options:**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `module` | String (Autocomplete) | ✅ | Select module to edit |
| `title` | String | ❌ | New title |
| `description` | String | ❌ | New description |
| `banner` | Attachment | ❌ | New banner image |

---

### `/module lesson edit`

**Purpose:** Edit an existing lesson.

**Options:**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `module` | String (Autocomplete) | ✅ | Select module |
| `lesson` | String (Autocomplete) | ✅ | Select lesson to edit |
| `title` | String | ❌ | New title |
| `content` | String | ❌ | New content |
| `image` | Attachment | ❌ | New image |

---

### `/module lesson reorder`

**Purpose:** Change the order of lessons (important for free vs premium access).

**Options:**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `module` | String (Autocomplete) | ✅ | Select module |
| `lesson` | String (Autocomplete) | ✅ | Select lesson to move |
| `new_position` | Integer | ✅ | New position (1 = free lesson) |

**⚠️ IMPORTANT:** Moving a lesson to position 1 makes it the FREE lesson!

---

### `/module list`

**Purpose:** List all modules with their status.

**Output:**
```
┌─────────────────────────────────────────────────────────────────────────────┐
│ 📚 **All Modules**                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ **Published:**                                                               │
│ 1. ✅ Trading Basics 101 (`trading-basics-101`) — 5 lessons                  │
│ 2. ✅ Technical Analysis (`technical-analysis`) — 8 lessons                  │
│                                                                              │
│ **Drafts:**                                                                  │
│ 3. 📝 Options Trading (`options-trading`) — 2 lessons                        │
│ 4. 📝 Crypto Basics (`crypto-basics`) — 0 lessons                            │
│                                                                              │
│ **Archived:**                                                                │
│ 5. 🗄️ Old Module (`old-module`)                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### `/module delete`

**Purpose:** Permanently delete a module.

**Options:**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `module` | String (Autocomplete) | ✅ | Select module to delete |
| `confirm` | Boolean | ✅ | Must be `true` to confirm |

**Flow:**
- If published → Delete forum post first
- Delete all lessons (CASCADE)
- Delete module record

---

## 🔄 FLOWCHARTS

### User Selects Lesson from Dropdown

```
USER CLICKS DROPDOWN → SELECTS LESSON
      │
      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         INTERACTION HANDLER                                  │
│                   CustomID: "Freemium:SelectLesson"                         │
└─────────────────────────────────────────────────────────────────────────────┘
      │
      │  Parse: module_id, lesson_id from interaction values
      │
      ▼
      │  STEP 1: Get User's Roles
      │
      │  member := s.GuildMember(guild_id, user_id)
      │  roles := member.Roles
      │
      ▼
      │  STEP 2: Determine Access Level
      │
      │  ┌─────────────────────────────────────────────┐
      │  │ Has LIFETIME (718643316786462772)?          │
      │  │ ├── YES → access_level = "lifetime"        │
      │  │ └── NO ↓                                   │
      │  │                                            │
      │  │ Has PREMIUM (885910828086362132)?          │
      │  │ ├── YES → access_level = "premium"         │
      │  │ └── NO → access_level = "free"             │
      │  └─────────────────────────────────────────────┘
      │
      ▼
      │  STEP 3: Get Lesson Data
      │
      │  SELECT * FROM freemium_lessons
      │  WHERE module_id = ? AND lesson_id = ?
      │
      ▼
      │  STEP 4: Check Access Permission
      │
      │  ┌─────────────────────────────────────────────────────────────────┐
      │  │                                                                 │
      │  │  IF access_level IN ("lifetime", "premium"):                   │
      │  │      → GRANT ACCESS (any lesson)                               │
      │  │                                                                 │
      │  │  ELSE IF lesson.order_index == 1:                              │
      │  │      → GRANT ACCESS (free lesson)                              │
      │  │                                                                 │
      │  │  ELSE:                                                         │
      │  │      → DENY ACCESS (show upgrade prompt)                       │
      │  │                                                                 │
      │  └─────────────────────────────────────────────────────────────────┘
      │
      ├─────────────────────────────────────────┐
      │                                         │
      ▼ ACCESS GRANTED                          ▼ ACCESS DENIED
      │                                         │
┌─────────────────────────────┐     ┌─────────────────────────────┐
│  SHOW LESSON CONTENT        │     │  SHOW UPGRADE PROMPT        │
│  (Ephemeral Message)        │     │  (Ephemeral Message)        │
│                             │     │                             │
│  ┌───────────────────────┐  │     │  ┌───────────────────────┐  │
│  │ [lesson_image]        │  │     │  │ 🔒 **Premium Only**   │  │
│  │                       │  │     │  │                       │  │
│  │ ## Lesson Title       │  │     │  │ This lesson requires  │  │
│  │                       │  │     │  │ a Premium membership. │  │
│  │ Full lesson content   │  │     │  │                       │  │
│  │ in markdown...        │  │     │  │ [🚀 Upgrade Now]      │  │
│  │                       │  │     │  │ (Link Button)         │  │
│  └───────────────────────┘  │     │  └───────────────────────┘  │
│                             │     │                             │
└─────────────────────────────┘     └─────────────────────────────┘
      │                                         │
      ▼                                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                         LOG ACCESS                               │
│                                                                  │
│  INSERT INTO freemium_access_log (                              │
│      user_id, module_id, lesson_id,                             │
│      access_granted, user_role_type                             │
│  ) VALUES (...)                                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
      │
      └─────▶ END
```

---

## 🎨 COMPONENTS V2 LAYOUTS

### Published Module Embed (Forum Post)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMPONENTS V2 MESSAGE STRUCTURE                           │
└─────────────────────────────────────────────────────────────────────────────┘

Container (type: 17, accent_color: brand_color)
│
├── MediaGallery (type: 12)
│   └── media_gallery_items: [
│         { media: { url: "attachment://freemium-banner.png" } }
│       ]
│
├── Separator (type: 14, spacing: "small")
│
├── TextDisplay (type: 10)
│   └── content: "### {MODULE_TITLE}\n\n{MODULE_DESCRIPTION}"
│
├── Separator (type: 14, spacing: "small")
│
├── TextDisplay (type: 10)
│   └── content: "**Select a lesson below:**\n✅ = Free | 🔒 = Premium"
│
├── ActionRow (type: 1)
│   └── StringSelect (type: 3)
│       ├── custom_id: "Freemium:SelectLesson:{module_id}"
│       ├── placeholder: "Choose a lesson..."
│       ├── min_values: 1
│       ├── max_values: 1
│       └── options: [
│             {
│               label: "1. Introduction",
│               value: "module_id:lesson_id",
│               description: "Learn the basics...",
│               emoji: { name: "✅" }   // FREE
│             },
│             {
│               label: "2. Understanding Charts",
│               value: "module_id:lesson_id",
│               description: "Dive into chart analysis...",
│               emoji: { name: "🔒" }   // PREMIUM
│             },
│             {
│               label: "3. Risk Management",
│               value: "module_id:lesson_id",
│               description: "Protect your capital...",
│               emoji: { name: "🔒" }   // PREMIUM
│             },
│             // ... more lessons
│           ]
│
├── Separator (type: 14, spacing: "small")
│
└── MediaGallery (type: 12)
    └── media_gallery_items: [
          { media: { url: "attachment://footer.png" } }
        ]
```

### Visual Representation

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                     [freemium-banner.png]                              │ │
│  │                                                                        │ │
│  │                     FREEMIUM MODULE BANNER                             │ │
│  │                                                                        │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  ### Trading Basics 101                                                     │
│                                                                              │
│  Master the fundamentals of trading with this comprehensive module.         │
│  Learn everything from basic chart reading to entry strategies.             │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  **Select a lesson below:**                                                 │
│  ✅ = Free | 🔒 = Premium                                                   │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  Choose a lesson...                                               ▼   │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  DROPDOWN EXPANDED:                                                         │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  ✅ 1. Introduction                                                    │ │
│  │     Learn the basics of trading                                        │ │
│  ├────────────────────────────────────────────────────────────────────────┤ │
│  │  🔒 2. Understanding Charts                                            │ │
│  │     Dive into chart analysis                                           │ │
│  ├────────────────────────────────────────────────────────────────────────┤ │
│  │  🔒 3. Risk Management                                                 │ │
│  │     Protect your capital                                               │ │
│  ├────────────────────────────────────────────────────────────────────────┤ │
│  │  🔒 4. Entry Strategies                                                │ │
│  │     When and how to enter                                              │ │
│  ├────────────────────────────────────────────────────────────────────────┤ │
│  │  🔒 5. Advanced Patterns                                               │ │
│  │     Complex pattern recognition                                        │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                        [footer.png]                                    │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Lesson Content Response (Ephemeral)

**For Granted Access:**
```
┌─────────────────────────────────────────────────────────────────────────────┐
│  EPHEMERAL MESSAGE — Only you can see this                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                     [lesson_image.png] (if exists)                     │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ## 📖 Introduction                                                         │
│                                                                              │
│  Welcome to Trading Basics 101! In this lesson, you'll learn...            │
│                                                                              │
│  **Key Points:**                                                            │
│  • Understanding market basics                                              │
│  • Types of trading                                                         │
│  • Setting up your first trade                                              │
│                                                                              │
│  [Full lesson content in markdown...]                                       │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│  📚 Module: Trading Basics 101 | Lesson 1 of 5                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**For Denied Access:**
```
┌─────────────────────────────────────────────────────────────────────────────┐
│  EPHEMERAL MESSAGE — Only you can see this                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  🔒 **Premium Content**                                                     │
│                                                                              │
│  **"Understanding Charts"** requires a Premium membership.                  │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  ✨ **Premium Benefits:**                                                   │
│  • Access ALL lessons in every module                                       │
│  • Exclusive trading signals                                                │
│  • Priority support                                                         │
│  • And much more!                                                           │
│                                                                              │
│  ┌─────────────────────┐                                                    │
│  │  🚀 Upgrade Now     │  ← Link Button to upgrade URL                      │
│  └─────────────────────┘                                                    │
│                                                                              │
│  💡 Tip: The first lesson of every module is FREE!                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## ✅ SETUP CHECKLIST

### Before You Start

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          SETUP CHECKLIST                                     │
└─────────────────────────────────────────────────────────────────────────────┘

□ DISCORD SETUP
  │
  ├── □ Forum Channel Created
  │   └── ID: 1447780920940695655
  │   └── Permissions: Bot can create threads, send messages
  │
  ├── □ Roles Verified
  │   ├── Lifetime Premium: 718643316786462772
  │   ├── Premium: 885910828086362132
  │   └── Free/Default: 718643370301325404
  │
  └── □ Bot Permissions
      ├── □ Manage Channels
      ├── □ Send Messages
      ├── □ Send Messages in Threads
      ├── □ Create Public Threads
      ├── □ Embed Links
      ├── □ Attach Files
      └── □ Read Message History

□ DATABASE SETUP
  │
  ├── □ PostgreSQL running
  ├── □ Tables created (run migrations)
  └── □ Default settings inserted

□ ASSETS PREPARED
  │
  ├── □ freemium-banner.png (default banner)
  │   └── Recommended: 1200x400px
  │
  └── □ footer.png (standard footer)
      └── Recommended: 1200x100px

□ ENVIRONMENT VARIABLES
  │
  ├── □ FREEMIUM_FORUM_CHANNEL_ID=1447780920940695655
  ├── □ LIFETIME_ROLE_ID=718643316786462772
  ├── □ PREMIUM_ROLE_ID=885910828086362132
  ├── □ FREE_ROLE_ID=718643370301325404
  └── □ UPGRADE_URL=https://your-upgrade-link.com
```

---

## 🖼️ ASSETS REQUIRED

### Image Assets

| Asset | Filename | Dimensions | Purpose |
|-------|----------|------------|---------|
| Default Banner | `freemium-banner.png` | 1200×400px | Default header for modules |
| Footer | `footer.png` | 1200×100px | Standard footer for all embeds |
| Upgrade CTA | `upgrade-cta.png` | 800×200px | Optional upgrade prompt image |

### Asset Placement

```
services/discordbot/assets/
├── freemium/
│   ├── freemium-banner.png
│   ├── footer.png
│   └── upgrade-cta.png
└── embed.go  (asset embedding)
```

---

## 📐 INTERACTION ID REFERENCE

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                         CUSTOM IDs — FREEMIUM SYSTEM                          │
└──────────────────────────────────────────────────────────────────────────────┘

LESSON SELECTION
└── Freemium:SelectLesson:{module_id} ──── Dropdown to select lesson

ADMIN MODALS (if needed for larger content)
├── Freemium:CreateModuleModal ─────────── Create module modal
├── Freemium:EditModuleModal:{module_id} ─ Edit module modal
├── Freemium:AddLessonModal:{module_id} ── Add lesson modal
└── Freemium:EditLessonModal:{module_id}:{lesson_id} ── Edit lesson

CONFIRMATION BUTTONS
├── Freemium:ConfirmDelete:{module_id} ─── Confirm module deletion
└── Freemium:CancelDelete ─────────────── Cancel deletion
```

---

## 🚀 QUICK START GUIDE

### Step 1: Create Your First Module

```bash
/module create name:"Trading Basics 101" description:"Learn the fundamentals of trading in this beginner-friendly module."
```

### Step 2: Add Lessons

```bash
# First lesson (FREE for everyone)
/module lesson add module:"trading-basics-101" title:"Introduction" content:"Welcome to Trading Basics..."

# Additional lessons (PREMIUM only)
/module lesson add module:"trading-basics-101" title:"Chart Reading" content:"Understanding candlesticks..."
/module lesson add module:"trading-basics-101" title:"Risk Management" content:"Never risk more than..."
```

### Step 3: Preview (Optional)

```bash
/module preview module:"trading-basics-101"
```

### Step 4: Publish

```bash
/module publish module:"trading-basics-101"
```

### Step 5: Verify

- Check the forum channel for the new post
- Test as a free member (should only access lesson 1)
- Test as premium member (should access all lessons)

---

## 📊 ANALYTICS QUERIES

### Most Popular Modules
```sql
SELECT 
    m.title,
    COUNT(*) as total_views,
    COUNT(CASE WHEN access_granted THEN 1 END) as successful_views
FROM freemium_access_log l
JOIN freemium_modules m ON l.module_id = m.module_id
GROUP BY m.module_id, m.title
ORDER BY total_views DESC;
```

### Conversion Funnel (Free → Premium)
```sql
SELECT 
    user_id,
    COUNT(CASE WHEN access_granted = FALSE THEN 1 END) as blocked_attempts,
    MIN(accessed_at) as first_blocked,
    MAX(accessed_at) as last_blocked
FROM freemium_access_log
WHERE access_granted = FALSE
GROUP BY user_id
ORDER BY blocked_attempts DESC
LIMIT 20;
```

### Lesson Popularity
```sql
SELECT 
    m.title as module,
    l.title as lesson,
    l.order_index,
    COUNT(*) as views
FROM freemium_access_log a
JOIN freemium_lessons l ON a.lesson_id = l.lesson_id AND a.module_id = l.module_id
JOIN freemium_modules m ON a.module_id = m.module_id
WHERE a.access_granted = TRUE
GROUP BY m.title, l.title, l.order_index
ORDER BY views DESC;
```

---

> **End of Freemium System Documentation**
> 
> This system allows you to build educational modules where free members
> can preview all content but only access the first lesson, encouraging
> upgrades to premium membership.





