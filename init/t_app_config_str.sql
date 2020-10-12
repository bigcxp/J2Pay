-- -------------------------------------------------------------
-- TablePlus 3.6.2(323)
--
-- https://tableplus.com/
--
-- Database: dc-wallet
-- Generation Time: 2020-06-24 11:01:59.7020
-- -------------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


INSERT INTO `t_app_config_str` (`id`, `k`, `v`) VALUES
('3', 'cold_wallet_address', '0xba71da7d56322dba3348204735bd419de245ad04'),
('4', 'hot_wallet_address', '0x47284d23b4c375878a52e7a7e5f0d4fbfb60fe22'),
('5', 'fee_wallet_address', '0xd3ecd590bd35c49733ae2cb82b7ae3f6bcec9b8c'),
('6', 'fee_wallet_address_list', '0xd3ecd590bd35c49733ae2cb82b7ae3f6bcec9b8c,0xd3ecd590bd35c49733ae2cb82b7ae3f6bcec9b8c'),
('7', 'cold_wallet_address_btc', 'n4LEjbHb2nFwKKrackfgoVCjpDKgDG9DZh'),
('8', 'hot_wallet_address_btc', 'mqzW5o46yBj7vwz4HyZNe7PYL25FHisvG4');


/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;