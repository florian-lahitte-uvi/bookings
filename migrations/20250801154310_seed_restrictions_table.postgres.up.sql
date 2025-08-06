-- Insert restrictions data only if they don't already exist
INSERT INTO public.restrictions (restriction_name,created_at,updated_at) 
SELECT 'Reservation','2025-07-31 00:00:00','2025-07-31 00:00:00'
WHERE NOT EXISTS (SELECT 1 FROM restrictions WHERE restriction_name = 'Reservation');

INSERT INTO public.restrictions (restriction_name,created_at,updated_at) 
SELECT 'Owner Block','2025-11-19 00:00:00','2025-11-19 00:00:00'
WHERE NOT EXISTS (SELECT 1 FROM restrictions WHERE restriction_name = 'Owner Block');
