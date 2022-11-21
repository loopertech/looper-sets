CREATE MIGRATION m1vvployqfrcm7wsbyuhkygs4lgqwcimosnyei7dokyldzc2f527aq
    ONTO m1t663bgubirb2qgh76uoqx2wdgnomrytf5suuzzdv4z62hv7txeya
{
  ALTER TYPE song::Song {
      ALTER LINK submitter {
          SET REQUIRED USING (SELECT
              user::User 
          LIMIT
              1
          );
      };
  };
};
