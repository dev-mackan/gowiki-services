
BEGIN TRANSACTION;

-- Step 1: Insert into the Text table
INSERT INTO Text (content) VALUES ('This is the initial content of the page.');

-- Step 2: Retrieve the last inserted text_id
-- Note: You will use this directly in the insert statements without using variables

-- Step 3: Insert into the Page table
INSERT INTO Page (title, latest_rev) VALUES ('My_First_Page', 0);

-- Step 4: Insert the initial revision into the Revision table
INSERT INTO Revision (page_id, text_id) 
VALUES (
    (SELECT last_insert_rowid()),  -- Get the last inserted page_id
    (SELECT last_insert_rowid() FROM Text ORDER BY text_id DESC LIMIT 1)  -- Get the last inserted text_id
);

-- Step 5: Update the latest_rev in the Page table
UPDATE Page 
SET latest_rev = (SELECT last_insert_rowid()) 
WHERE page_id = (SELECT last_insert_rowid());

COMMIT;
