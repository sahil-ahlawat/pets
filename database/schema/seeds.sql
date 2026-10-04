-- seeds.sql
-- Initial seed data for Pet Digital Identity Platform

-- Standard Core Canine Vaccination Types
INSERT INTO vaccination_types (species, name, code, is_core, recommended_interval_months, description)
VALUES 
    ('DOG', 'Rabies (Canine 1-Year/3-Year)', 'RABIES_CANINE', TRUE, 36, 'Mandatory rabies vaccination protecting against the rabies virus.'),
    ('DOG', 'DHPP (Distemper, Hepatitis, Parvovirus, Parainfluenza)', 'DHPP', TRUE, 12, 'Core 5-in-1 combo vaccine guarding puppies and dogs against fatal viral infections.'),
    ('DOG', 'Bordetella (Kennel Cough)', 'BORDETELLA', FALSE, 6, 'Recommended for dogs frequenting boarding, daycare, or parks.'),
    ('DOG', 'Leptospirosis', 'LEPTO', FALSE, 12, 'Bacterial infection protection, especially for outdoor and water-loving dogs.'),
    ('DOG', 'Lyme Disease', 'LYME', FALSE, 12, 'Protects against tick-borne Borrelia burgdorferi.')
ON CONFLICT (species, code) DO NOTHING;

-- Standard Core Feline Vaccination Types
INSERT INTO vaccination_types (species, name, code, is_core, recommended_interval_months, description)
VALUES 
    ('CAT', 'Rabies (Feline Purevax/Standard)', 'RABIES_FELINE', TRUE, 12, 'Crucial rabies prophylaxis for domestic and indoor/outdoor cats.'),
    ('CAT', 'FVRCP (Feline Viral Rhinotracheitis, Calicivirus, Panleukopenia)', 'FVRCP', TRUE, 36, 'Core 3-in-1 feline respiratory and distemper protection.'),
    ('CAT', 'FeLV (Feline Leukemia Virus)', 'FELV', FALSE, 12, 'Highly recommended for cats with outdoor access or multiple cat households.')
ON CONFLICT (species, code) DO NOTHING;

-- Seed Sample Veterinary Clinics
INSERT INTO vets (
    clinic_name, doctor_name, license_number, phone_number, emergency_phone,
    email, website_url, address_line1, city, state, postal_code, country_code,
    latitude, longitude, is_24h_emergency, specializations, rating, review_count, is_verified
) VALUES 
(
    'Metropolitan Animal Emergency & Specialty Hospital',
    'Dr. Elena Rostova, DVM, DACVECC',
    'VET-883921',
    '+1 (555) 234-5678',
    '+1 (555) 234-9999',
    'care@metrovet-emergency.com',
    'https://metrovet-emergency.com',
    '742 Evergreen Promenade, Suite 100',
    'San Francisco',
    'CA',
    '94107',
    'US',
    37.7749,
    -122.4194,
    TRUE,
    ARRAY['Emergency Medicine', 'Critical Care', 'Orthopedic Surgery', 'Ultrasound'],
    4.95,
    148,
    TRUE
),
(
    'Paws & Claws Holistic Veterinary Center',
    'Dr. Marcus Vance, MRCVS',
    'VET-419082',
    '+1 (555) 789-0123',
    '+1 (555) 789-0124',
    'hello@pawsandclaws-vet.com',
    'https://pawsandclaws-vet.com',
    '1248 Oak Boulevard',
    'San Francisco',
    'CA',
    '94110',
    'US',
    37.7599,
    -122.4148,
    FALSE,
    ARRAY['General Practice', 'Dental Surgery', 'Feline Medicine', 'Preventative Care'],
    4.88,
    92,
    TRUE
)
ON CONFLICT DO NOTHING;

-- Seed Blog Categories
INSERT INTO blog_categories (slug, name, description, icon_name, display_order)
VALUES 
    ('pet-health', 'Pet Health & Wellness', 'Expert veterinary guidance on vaccinations, preventive care, and wellness.', 'heart-pulse', 1),
    ('lost-and-found', 'Lost Pet Recovery', 'Action plans, search strategies, and recovery stories to keep pets safe.', 'shield-alert', 2),
    ('smart-tags', 'Smart QR & NFC Technology', 'How modern digital pet identity outperforms obsolete metal collar tags.', 'qr-code', 3),
    ('nutrition-lifestyle', 'Nutrition & Lifestyle', 'Dietary recommendations, exercise, and enriching your pet companionship.', 'bone', 4)
ON CONFLICT (slug) DO NOTHING;

-- Seed Sample Blog Post
INSERT INTO blog_posts (
    title, slug, excerpt, content_markdown, reading_time_minutes, is_published, published_at, views_count
) VALUES (
    'Why Metal Tags Fail: The Digital QR Revolution in Pet Safety',
    'why-metal-tags-fail-digital-qr-revolution',
    'Traditional metal tags scratch, wear off, and reveal only a static phone number. Here is how modern QR digital identity instantly re-unites lost pets.',
    '# Why Metal Tags Fail: The Digital QR Revolution in Pet Safety

Every year, millions of family pets wander away or slip out through an open gate. In past decades, their only line of defense was an engraved metal disc stamped with a phone number that quickly becomes scratched and illegible.

### The Limitations of Legacy Tags
* **Static Information**: When you travel, move houses, or change phone numbers, metal tags are immediately obsolete.
* **Zero Medical Context**: If your pet is diabetic or has severe medication allergies, a regular tag communicates nothing to the rescuer.
* **Privacy Risks**: Displaying your private home address and mobile number on an open collar invites privacy and stalking concerns.

### The Power of Dynamic Digital QR Identity
With dynamic QR tags, scanning the collar opens a **secure, real-time responsive profile**. The owner can update contact details, activate **Lost Mode**, and receive instant scan alerts with approximate location coordinates.',
    4,
    TRUE,
    CURRENT_TIMESTAMP,
    342
)
ON CONFLICT (slug) DO NOTHING;
