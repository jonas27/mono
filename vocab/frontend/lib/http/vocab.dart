import 'package:grpc/grpc_web.dart';
import '../proto/vocab.pbgrpc.dart';

Future<void> getVocabList() async {
  Uri uri = Uri.parse('http://localhost:8082');
  final channel = GrpcWebClientChannel.xhr(uri);
  final client = VocabServiceClient(channel);

  final request = VocabListRequest();
  final stream = client.listVocabs(request);

  await for (final message in stream) {
    // Handle message received from stream
    print('Received message: \n$message');
  }
  await channel.shutdown();
}
