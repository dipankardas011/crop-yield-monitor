from flask import Flask, request, render_template
import os

app = Flask(__name__)

UPLOAD_FOLDER = 'uploads'
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

# Counter to keep track of uploaded files
file_counter = 1

# Function to generate a default filename with serial number
def generate_filename():
    global file_counter
    filename = f"uploaded_image_{file_counter}.png"  # You can change the file extension as needed
    file_counter += 1
    return filename

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/upload', methods=['POST'])
def upload_file():
    if 'image' not in request.files:
        return 'No file part'

    file = request.files['image']

    if file.filename == '':
        return 'No selected file'

    if file:
        # Create 'uploads' directory if it doesn't exist
        if not os.path.exists(app.config['UPLOAD_FOLDER']):
            os.makedirs(app.config['UPLOAD_FOLDER'])

        filename = os.path.join(app.config['UPLOAD_FOLDER'], generate_filename())
        file.save(filename)
        return 'File uploaded successfully'

if __name__ == '__main__':
    app.run(debug=True)
