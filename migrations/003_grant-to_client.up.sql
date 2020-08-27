\c msg_service;

BEGIN;

GRANT ALL ON TABLE public.users TO service_client;
GRANT ALL ON TABLE public.chat TO service_client;
GRANT ALL ON TABLE public.chat_users TO service_client;
GRANT ALL ON TABLE public.msg TO service_client;

COMMIT;