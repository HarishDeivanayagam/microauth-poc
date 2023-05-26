ALTER TABLE member_invite ADD COLUMN organization_id VARCHAR(255) NOT NULL;
ALTER TABLE member_invite ADD CONSTRAINT fk_organization_id FOREIGN KEY (organization_id) REFERENCES organizations (id);
