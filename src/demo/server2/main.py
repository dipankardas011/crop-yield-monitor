from flask import Flask, request, jsonify
from flask_cors import CORS  # Import the CORS extension

app = Flask(__name__)
CORS(app, resources={r"/upload": {"origins": "*"}})  # Enable CORS for all routes in the app

@app.route('/upload', methods=['POST'])
def upload_image():
    # Get headers
    auth_header = request.headers.get('Authorization')
    
    # Get query parameters
    username_param = request.args.get('username')
    
    print("username= ",username_param)
    print("auth_header= ", auth_header)

    data = request.get_json()
    if not data or 'image' not in data or 'format' not in data:
        return jsonify({'error': 'Invalid request'}), 400

    image_data = bytes(data['image'])
    image_format = data['format']

    file_extension = 'jpeg' if image_format == 'image/jpeg' else 'png'
    file_name = f'image.{file_extension}'

    try:
        with open(file_name, 'wb') as image_file:
            image_file.write(image_data)
    except Exception as e:
        print(e)
        return jsonify({'error': 'Failed to save image'}), 500

    return jsonify({'status': 'Image uploaded', 'crops': ['crop', 'maize']}), 200

if __name__ == '__main__':
    app.run(debug=True)
