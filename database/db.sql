-- Use the beanbag database
\c beanbag

-- QUIZZES TABLE
CREATE TABLE quizzes (
    quiz_id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) UNIQUE NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    is_priv BOOLEAN NOT NULL DEFAULT true,
    timer INTEGER NOT NULL DEFAULT 30
);

-- QUESTIONS TABLE (Depends on QUIZZES)
CREATE TABLE questions (
    ques_id SERIAL PRIMARY KEY,
    quiz_id INTEGER REFERENCES quizzes(quid_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    sp_timer_option BOOLEAN NOT NULL DEFAULT false,
    sp_timer INTEGER NOT NULL DEFAULT 30
);

-- ANWSERS TABLE (Depends on QUESTIONS)
CREATE TABLE answers (
    ans_id SERIAL PRIMARY KEY,
    ques_id INTEGER REFERENCES questions(quies_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL
);

-- -- ORDER_ITEMS TABLE (linking items with orders)
-- CREATE TABLE order_items (
--     order_item_id SERIAL PRIMARY KEY,
--     order_id INTEGER REFERENCES orders(order_id),
--     item_id INTEGER REFERENCES items(item_id),
--     quantity INTEGER NOT NULL,
--     item_price DECIMAL(10, 2) NOT NULL
-- );

-- Set correct ownership of tables to postgres
ALTER TABLE quizzes OWNER TO postgres;
ALTER TABLE questions OWNER TO postgres;
ALTER TABLE answers OWNER TO postgres;
-- ALTER TABLE order_items OWNER TO garrick;

-- Insert sample magical items
-- INSERT INTO items (image, name, description, price, stock)
-- VALUES
-- ('9.png', 'Sword of Eternal Frost',
--  'A blade forged from the glaciers of the North. Its edge never dulls and can chill foes to the bone.',
--  200.00,
--  3),
-- ('8.png','Amulet of the Phoenix',
--  'Glowing with a fiery radiance, it grants its bearer resistance to fire and a hint of rebirth.',
--  120.00,
--  5),
-- ('10.png','Elixir of Arcane Might',
--  'A swirling potion that enhances magical ability for a short time, at a risk of arcane backlash.',
--  40.00,
--  10),
-- ('1.png','Ancient Dragon Scale Shield',
--  'Forged from the scales of the Red Highlands dragon. Highly resistant to flame and claw.',
--  300.00,
--  2),
-- ('2.png','Ring of Whispering Winds',
--  'Carries the voices of distant lands to the wearer, and sometimes, voices from beyond.',
--  80.00,
--  8),
-- ('3.png','Shadows Embrace Cloak',
--  'A cloak woven from the threads of night, allowing the wearer to blend into darkness and silence their footsteps.',
--  150.00,
--  4),
-- ('4.png','Emerald Eye of Insight',
--  'A brilliant emerald orb granting short bursts of enhanced perception, revealing hidden truths and illusions.',
--  95.00,
--  6),
-- ('5.png','Golemheart Gauntlets',
--  'Sturdy gloves embedded with enchanted stone shards, boosting the wearers physical strength and durability.',
--  110.00,
--  5),
-- ('6.png','Hammer of Thunderous Wrath',
--  'A heavy warhammer that crackles with bound lightning. Strikes reverberate like thunder, staggering opponents.',
--  220.00,
--  3),
-- ('7.png','Wand of Illusory Dreams',
--  'A slender wand that conjures illusions from the wielders imagination, though at risk of blending dream and reality.',
--  75.00,
--  7);



-- Ensure privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO postgres;

-- Explicitly grant privileges to the user
grant all privileges on all tables in schema public to postgres;
grant all privileges on all sequences in schema public to postgres;
