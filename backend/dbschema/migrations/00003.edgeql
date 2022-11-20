CREATE MIGRATION m1t663bgubirb2qgh76uoqx2wdgnomrytf5suuzzdv4z62hv7txeya
    ONTO m1ggwu5bc3l5wi6j5km3o3nfs62bums3lbbuhp33ehyklmgeq43qjq
{
  ALTER TYPE user::User {
      ALTER PROPERTY display_name {
          DROP CONSTRAINT std::max_len_value(3);
      };
  };
  ALTER TYPE user::User {
      ALTER PROPERTY display_name {
          CREATE CONSTRAINT std::max_len_value(32);
      };
  };
  ALTER TYPE user::User {
      ALTER PROPERTY username {
          DROP CONSTRAINT std::max_len_value(3);
      };
  };
  ALTER TYPE user::User {
      ALTER PROPERTY username {
          CREATE CONSTRAINT std::max_len_value(32);
      };
  };
};
