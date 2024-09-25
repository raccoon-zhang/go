//
//  prefixView.swift
//  quickPaste
//
//  Created by zhangshaohua on 2024/9/26.
//

import Foundation
import Cocoa
import SnapKit

class PrefixView: NSView {
    
    public var textChangeBlock: ((String) -> Void)?
    
    private lazy var label = {
        let label = Label()
        label.stringValue = "Prefix:"
        return label
    }()
    
    private lazy var inputView = {
        let view = NSTextField()
        view.wantsLayer = true
        view.layer?.borderColor = CGColor.white
        view.layer?.borderWidth = 1.0
        view.delegate = self
        return view
    }()
    
    init() {
        super.init(frame: NSRect(x: 0, y: 0, width: 200, height: 20))
        self.setupUI()
    }
    
    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
    
    public func setupInitialValue(_ value: String) {
        inputView.stringValue = value
    }
    
    private func setupUI() {
        addSubview(label)
        addSubview(inputView)
        label.snp.makeConstraints {
            $0.left.equalToSuperview().offset(15)
            $0.centerY.equalToSuperview()
        }
        inputView.snp.makeConstraints {
            $0.left.equalTo(label.snp.right).offset(10)
            $0.centerY.right.equalToSuperview()
        }
    }
}

extension PrefixView: NSTextFieldDelegate {
    func controlTextDidEndEditing(_ obj: Notification) {
        self.textChangeBlock?(inputView.stringValue)
    }
}
