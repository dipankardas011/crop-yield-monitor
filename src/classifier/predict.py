
# Load libraries
from PIL import Image
import torch
import glob
from torchvision.transforms import transforms
import torch.optim
from torch.autograd import Variable
import pathlib

from src.ConvNet import ConvNet

batchSize = 4

train_path = './references/train'

root = pathlib.Path(train_path)
classes = sorted([j.name.split('/')[-1] for j in root.iterdir()])

# CNN Network
print(classes)
checkpoint = torch.load('models/best_checkpoint.model')
model = ConvNet(num_classes=len(classes), batchsize=batchSize)
model.load_state_dict(checkpoint)
model.eval()

transformer_pred = transforms.Compose([
    transforms.Resize((150, 150)),
    transforms.ToTensor(),  # 0-255 to 0-1, numpy to tensors
    transforms.Normalize([0.5, 0.5, 0.5],  # 0-1 to [-1,1] , formula (x-mean)/std
                         [0.5, 0.5, 0.5])
])
pred_path = './references/pred/'


# prediction function

def prediction(img_path, transformer):
    image = Image.open(img_path)
    image_tensor = transformer(image).float()
    image_tensor = image_tensor.unsqueeze_(0)
    if torch.cuda.is_available():
        image_tensor.cuda()
    input = Variable(image_tensor)
    output = model(input)
    # print(f"Output Array {output.data.numpy()}")
    index = output.data.numpy().argmax()
    predClass = classes[index]

    # print(f"Input Img: {img_path[img_path.rfind('/') + 1:]}\nPredicted class {predClass}\nindex {index}",
    #       end="\n=========\n")
    return predClass


images_path = glob.glob(pred_path + '/*.jpg')

pred_accuracy = 0.0

for i in images_path:
    predClass = prediction(i, transformer_pred)
    imgName = i[i.rfind('/') + 1:]
    orgClass = imgName.split('_')[:1][0].title() + " Soil"

    if orgClass == predClass:
        pred_accuracy += 1

print(f"Correct Prediction: {pred_accuracy} out of {len(images_path)} -> {pred_accuracy/len(images_path)}")
