DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'contacts') THEN
        CREATE TABLE contacts (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            contact_username VARCHAR(255) NOT NULL,
            status VARCHAR(50) NOT NULL DEFAULT 'pending', /* 'pending', 'accepted', 'rejected', etc. */
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        );

        -- Ensure that a user can't have duplicate requests or contacts
        CREATE UNIQUE INDEX unique_contact_request
        ON contacts (username, contact_username);
    END IF;
END $$;
