CREATE MIGRATION m1ggwu5bc3l5wi6j5km3o3nfs62bums3lbbuhp33ehyklmgeq43qjq
    ONTO m1uax3mgynkmnenxmp2anm5altg4xzgafofae53zj3riz72oqvktmq
{
  ALTER TYPE song::Song {
      ALTER PROPERTY music_link {
          RENAME TO music_url;
      };
  };
  ALTER TYPE song::Song {
      ALTER PROPERTY song_link {
          RENAME TO song_url;
      };
  };
  ALTER TYPE song::Song {
      ALTER PROPERTY video_link {
          RENAME TO video_url;
      };
  };
  ALTER TYPE user::User {
      ALTER PROPERTY displayName {
          RENAME TO display_name;
      };
  };
};
