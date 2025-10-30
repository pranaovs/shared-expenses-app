import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class StorageService {
  static const _storage = FlutterSecureStorage();
  static const String _tokenKey = 'jwt_token';
  static const String _serverUrlKey = 'server_url';

  Future<void> saveToken(String token) async {
    await _storage.write(key: _tokenKey, value: token);
  }

  Future<String?> getToken() async {
    return await _storage.read(key: _tokenKey);
  }

  Future<void> deleteToken() async {
    await _storage.delete(key: _tokenKey);
  }

  Future<void> saveServerUrl(String url) async {
    await _storage.write(key: _serverUrlKey, value: url);
  }

  Future<String?> getServerUrl() async {
    return await _storage.read(key: _serverUrlKey);
  }

  Future<void> clearAll() async {
    await _storage.deleteAll();
  }
}
