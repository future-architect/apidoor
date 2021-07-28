BEGIN;
INSERT INTO public.apiinfo
    (id, name, source, description, thumbnail)
VALUES
    (
        3,
        'Awesome API',
        'Nice Company',
        'provide fantastic information.',
        'test.com/img/123'
    ),
    (
        4,
        'Awesome API v2',
        'Nice Company',
        'provide special information.',
        'test.com/img/456'
    );
END;