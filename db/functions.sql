CREATE OR REPLACE FUNCTION get_users_by_department(dept_id UUID)
RETURNS TABLE(
    user_id UUID,
    first_name TEXT,
    last_name TEXT,
    department UUID,
    email TEXT,
    phone TEXT,
    user_pic UUID,
    properties JSONB,
    active_tickets BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.user_id,
        u.first_name,
        u.last_name,
        u.department,
        u.email,
        u.phone,
        u.user_pic,
        CASE 
            WHEN u.latest_ticket IS NULL THEN '{}'::jsonb
            ELSE jsonb_build_object(
                'status', t.status,
                'workstarted_at', t.workstarted_at,
                'workfinished_at', t.workfinished_at,
                'client_name', t.client_name
            )
        END AS properties,
        (
            SELECT COUNT(*) 
            FROM tickets 
            WHERE executor = u.user_id 
            AND status NOT IN ('closed', 'cancelled')
        ) AS active_tickets
    FROM 
        users u
    LEFT JOIN tickets t ON u.latest_ticket = t.id
    WHERE 
        u.department = dept_id;
END;
$$ LANGUAGE plpgsql;
