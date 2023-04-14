# Load libraries
from PIL import Image
import os
import numpy as np
import torch
import glob
import torch.nn as nn
from torchvision.transforms import transforms
from torch.utils.data import DataLoader
import torch.optim
from torch.autograd import Variable
import torchvision
import pathlib

# checking for device
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

# Transforms
transformer = transforms.Compose([
    transforms.Resize((150, 150)),
    transforms.ColorJitter(),
    transforms.RandomHorizontalFlip(),
    transforms.RandomVerticalFlip(),
    transforms.RandomRotation(+90),
    transforms.ToTensor(),  # 0-255 to 0-1, numpy to tensors
    transforms.Normalize([0.5, 0.5, 0.5],  # 0-1 to [-1,1] , formula (x-mean)/std
                         [0.5, 0.5, 0.5])
])

# Dataloader

# Path for training and testing directory
train_path = './references/train'
test_path = './references/test'

batchSize = 4

train_loader = DataLoader(
    torchvision.datasets.ImageFolder(train_path, transform=transformer),
    batch_size=batchSize, shuffle=True
)
test_loader = DataLoader(
    torchvision.datasets.ImageFolder(test_path, transform=transformer),
    batch_size=batchSize // 2, shuffle=True
)

# categories
root = pathlib.Path(train_path)
classes = sorted([j.name.split('/')[-1] for j in root.iterdir()])

# CNN Network
print(classes)


class ConvNet(nn.Module):
    def __init__(self, num_classes):
        super(ConvNet, self).__init__()

        # Output size after convolution filter
        # ((w-f+2P)/s) +1

        # Input shape= (256,3,150,150)

        self.conv1 = nn.Conv2d(
            in_channels=3, out_channels=12, kernel_size=3, stride=1, padding=1)
        # Shape= (256,12,150,150)
        self.bn1 = nn.BatchNorm2d(num_features=12)
        # Shape= (256,12,150,150)
        self.relu1 = nn.ReLU()
        # Shape= (256,12,150,150)

        self.pool = nn.MaxPool2d(kernel_size=2)
        # Reduce the image size be factor 2
        # Shape= (256,12,75,75)

        self.conv2 = nn.Conv2d(
            in_channels=12, out_channels=20, kernel_size=3, stride=1, padding=1)
        # Shape= (256,20,75,75)
        self.relu2 = nn.ReLU()
        # Shape= (256,20,75,75)

        self.conv3 = nn.Conv2d(
            in_channels=20, out_channels=batchSize, kernel_size=3, stride=1, padding=1)
        # Shape= (256,32,75,75)
        self.bn3 = nn.BatchNorm2d(num_features=batchSize)
        # Shape= (256,32,75,75)
        self.relu3 = nn.ReLU()
        # Shape= (256,32,75,75)

        self.fc = nn.Linear(in_features=75 * 75 * batchSize, out_features=num_classes)

        # softmax & sigmoid
        # Feed forward function

    def forward(self, input):
        output = self.conv1(input)
        output = self.bn1(output)
        output = self.relu1(output)

        output = self.pool(output)

        output = self.conv2(output)
        output = self.relu2(output)

        output = self.conv3(output)
        output = self.bn3(output)
        output = self.relu3(output)

        # Above output will be in matrix form, with shape (256,32,75,75)

        output = output.view(-1, batchSize * 75 * 75)

        output = self.fc(output)

        return output


model = ConvNet(num_classes=len(classes)).to(device)

# Optmizer and loss function
optimizer = torch.optim.Adam(model.parameters(), lr=0.001)
loss_function = nn.CrossEntropyLoss()

# calculating the size of training and testing images
train_count = len(glob.glob(train_path + '/**/*.jpg'))
test_count = len(glob.glob(test_path + '/**/*.jpg'))

print(train_count, test_count)

# Model training and saving best model

best_accuracy = 0.0
num_epochs = 100

for epoch in range(num_epochs):

    # Evaluation and training on training dataset
    model.train()
    train_accuracy = 0.0
    train_loss = 0.0

    for i, (images, labels) in enumerate(train_loader):
        if torch.cuda.is_available():
            images = Variable(images.cuda())
            labels = Variable(labels.cuda())

        optimizer.zero_grad()

        outputs = model(images)
        loss = loss_function(outputs, labels)
        loss.backward()
        optimizer.step()

        train_loss += loss.cpu().data * images.size(0)
        _, prediction = torch.max(outputs.data, 1)

        train_accuracy += int(torch.sum(prediction == labels.data))

    train_accuracy = train_accuracy / train_count
    train_loss = train_loss / train_count

    # Evaluation on testing dataset
    model.eval()

    test_accuracy = 0.0
    for i, (images, labels) in enumerate(test_loader):
        if torch.cuda.is_available():
            images = Variable(images.cuda())
            labels = Variable(labels.cuda())

        outputs = model(images)
        _, prediction = torch.max(outputs.data, 1)
        test_accuracy += int(torch.sum(prediction == labels.data))

    test_accuracy = test_accuracy / test_count

    print('Epoch: ' + str(epoch) + ' Train Loss: ' + str(train_loss) +
          ' Train Accuracy: ' + str(train_accuracy * 100) + ' Test Accuracy: ' + str(test_accuracy * 100))

    # Save the best model
    if test_accuracy > best_accuracy:
        torch.save(model.state_dict(), 'best_checkpoint.model')
        best_accuracy = test_accuracy

print("DONE WITH PREDICTION")
print(f"==================Best testAcc [{best_accuracy * 100} %]==================")

checkpoint = torch.load('best_checkpoint.model')
model = ConvNet(num_classes=len(classes))
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
    index = output.data.numpy().argmax()
    pred = classes[index]
    return pred


images_path = glob.glob(pred_path + '/*.jpg')

pred_dict = {}

for i in images_path:
    pred_dict[i[i.rfind('/') + 1:]] = prediction(i, transformer_pred)

print(pred_dict)
