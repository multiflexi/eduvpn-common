<?php declare(strict_types=1);

use EduVpn\Common\Discovery;
use EduVpn\Common\InvalidSignatureException;
use EduVpn\Common\InvalidSignatureUnknownKeyException;
use EduVpn\Common\SignatureTooOldException;
use PHPUnit\Framework\TestCase;

class DiscoveryTest extends TestCase
{
	/** @var string */
	private const TEST_DATA_DIR = '../../test_data';

	public static function setUpBeforeClass(): void
	{
		preg_match('/[\r\n](\S+)\s*/', file_get_contents(self::TEST_DATA_DIR . '/dummy/public.key'), $matches);
		Discovery::insecureTestingSetExtraKey($matches[1]);
	}

	public function testValid(): void
	{
		$this->expectNotToPerformAssertions();
		Discovery::verify(file_get_contents(self::TEST_DATA_DIR . '/dummy/server_list.json.minisig'),
			  file_get_contents(self::TEST_DATA_DIR . '/dummy/server_list.json'),
			  'server_list.json', 0);
	}

	public function testInvalidSignature(): void
	{
		$this->expectException(InvalidSignatureException::class);
		Discovery::verify(file_get_contents(self::TEST_DATA_DIR . '/dummy/random.txt'),
			  file_get_contents(self::TEST_DATA_DIR . '/dummy/server_list.json'),
			  'server_list.json', 0);
	}

	public function testWrongKey(): void
	{
		$this->expectException(InvalidSignatureUnknownKeyException::class);
		Discovery::verify(file_get_contents(self::TEST_DATA_DIR . '/dummy/server_list.json.wrong_key.minisig'),
			  file_get_contents(self::TEST_DATA_DIR . '/dummy/server_list.json'),
			  'server_list.json', 0);
	}

	#[deprecated]
	public function testOldSignature(): void
	{
		$this->expectException(SignatureTooOldException::class);
		Discovery::verify(file_get_contents(self::TEST_DATA_DIR . '/dummy/server_list.json.minisig'),
			  file_get_contents(self::TEST_DATA_DIR . '/dummy/server_list.json'),
			  'server_list.json', 1 << 31);
	}

	public function testUnknownExpectedFileName(): void
	{
		$this->expectException(InvalidArgumentException::class);
		Discovery::verify(file_get_contents(self::TEST_DATA_DIR . '/dummy/other_list.json.minisig'),
			  file_get_contents(self::TEST_DATA_DIR . '/dummy/other_list.json'),
			  'other_list.json', 0);
	}
}
