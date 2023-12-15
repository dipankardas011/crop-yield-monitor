from flask import Flask, request, jsonify
from flask_cors import CORS  # Import the CORS extension
from SoilRecommendation import predict_single_image, rf_classifier, label_encoder

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
    if not data or 'raw_image_bytes' not in data or 'image_format' not in data:
        return jsonify({'error': 'Invalid request'}), 400

    image_data = bytes(data['raw_image_bytes'])
    image_format = data['image_format']

    file_extension = 'jpeg' if image_format == 'image/jpeg' else 'png'
    file_name = f'image.{file_extension}'

    try:
        with open(file_name, 'wb') as image_file:
            image_file.write(image_data)
    except Exception as e:
        print(e)
        return jsonify({'error': 'Failed to save image'}), 500

    # Make a prediction for the test image
    predicted_class = predict_single_image("/app/"+file_name, rf_classifier, label_encoder)

    # Print the predicted class
    print(f"The test image belongs to the class: {predicted_class}")

    if predicted_class == "Gravel":
        return jsonify({
            'Crops': ['Grapevines', 'certain varieties of strawberries', 'tomatoes', 'peppers', 'rosemary'],
            'Message':'Gravelly soil tends to be well-draining but may not retain water and nutrients as effectively. Crops that can tolerate well-draining conditions and do not require high fertility may be suitable.',
            'Status': 'Ready'}), 200

    if predicted_class == "Sand":
        return jsonify({
            'Crops': ['Carrots', 'radishes', 'potatoes', 'asparagus', 'melons', 'lavender'],
            'Message':'Sandy soil drains quickly but may lack fertility and moisture retention. Crops that thrive in well-draining conditions and can tolerate lower fertility are suitable for sandy soil.',
            'Status': 'Ready'}), 200

    if predicted_class == "Silt":
        return jsonify({
            'Crops': ['Beans', 'peas', 'lettuce', 'spinach', 'broccoli', 'cabbage', 'most garden vegetables'],
            'Message':'Silt has finer particles and better moisture retention than sand. It is fertile and suitable for a wide range of crops.',
            'Status': 'Ready'}), 200

if __name__ == '__main__':
    app.run(debug=True)
